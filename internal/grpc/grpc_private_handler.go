package grpc

import (
	"context"
	"encoding/json"

	"github.com/flash1nho/GophKeeper/config"
	"github.com/flash1nho/GophKeeper/internal/facade"
	"github.com/flash1nho/GophKeeper/internal/models/secrets"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

var typeToSecret = map[string]func(userID int, masterKey []byte) secrets.Secret{
	"Text": func(u int, k []byte) secrets.Secret { return secrets.NewText(u, k) },
	"Cred": func(u int, k []byte) secrets.Secret { return secrets.NewCred(u, k) },
	"Card": func(u int, k []byte) secrets.Secret { return secrets.NewCard(u, k) },
}

type GrpcPrivateHandler struct {
	UnimplementedGophKeeperPrivateServiceServer

	pool     *pgxpool.Pool
	settings config.SettingsObject
	facade   *facade.Facade
}

func NewGrpcPrivateHandler(pool *pgxpool.Pool, settings config.SettingsObject, facade *facade.Facade) *GrpcPrivateHandler {
	return &GrpcPrivateHandler{
		pool:     pool,
		settings: settings,
		facade:   facade,
	}
}

func (g *GrpcPrivateHandler) Create(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
	secret, err := g.prepareSecret(ctx, req.Type, req.Data)

	if err != nil {
		return nil, err
	}

	res, err := secrets.Create(ctx, g.pool, secret)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка при создании секрета: %v", err)
	}

	protoData, err := g.mapToProto(res)

	if err != nil {
		return nil, err
	}

	return &CreateResponse{Secrets: protoData}, nil
}

func (g *GrpcPrivateHandler) Get(ctx context.Context, req *GetRequest) (*GetResponse, error) {
	secret, err := g.getSecretInstance(ctx, req.Type)

	if err != nil {
		return nil, err
	}

	res, err := secrets.Get(ctx, g.pool, secret, int(req.ID))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "не удалось получить секрет: %v", err)
	}

	protoData, err := g.mapToProto(res)

	if err != nil {
		return nil, err
	}

	return &GetResponse{Secrets: protoData}, nil
}

func (g *GrpcPrivateHandler) List(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	secret, err := g.getSecretInstance(ctx, req.Type)

	if err != nil {
		return nil, err
	}

	res, err := secrets.List(ctx, g.pool, secret)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "не удалось получить список секретов: %v", err)
	}

	protoData, err := g.mapToProto(res)

	if err != nil {
		return nil, err
	}

	return &ListResponse{Secrets: protoData}, nil
}

func (g *GrpcPrivateHandler) Update(ctx context.Context, req *UpdateRequest) (*UpdateResponse, error) {
	secret, err := g.prepareSecret(ctx, req.Type, req.Data)

	if err != nil {
		return nil, err
	}

	res, err := secrets.Update(ctx, g.pool, secret, int(req.ID))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка при обновлении секрета: %v", err)
	}

	protoData, err := g.mapToProto(res)
	if err != nil {
		return nil, err
	}

	return &UpdateResponse{Secrets: protoData}, nil
}

func (g *GrpcPrivateHandler) Delete(ctx context.Context, req *DeleteRequest) (*DeleteResponse, error) {
	secret, err := g.getSecretInstance(ctx, req.Type)

	if err != nil {
		return nil, err
	}

	if err := secrets.Delete(ctx, g.pool, secret, int(req.ID)); err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка при удалении секрета: %v", err)
	}

	return &DeleteResponse{}, nil
}

func (g *GrpcPrivateHandler) getSecretInstance(ctx context.Context, secretType string) (secrets.Secret, error) {
	userID, err := g.facade.GetUserIDFromContext(ctx)

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "пользователь не авторизован")
	}

	fn, ok := typeToSecret[secretType]

	if !ok {
		return nil, status.Error(codes.InvalidArgument, "неизвестный тип секрета")
	}

	return fn(userID, g.settings.MasterKey), nil
}

func (g *GrpcPrivateHandler) prepareSecret(ctx context.Context, secretType string, data *structpb.Struct) (secrets.Secret, error) {
	secret, err := g.getSecretInstance(ctx, secretType)
	g.mergeAdditionalData(secretType, data)

	if err != nil {
		return nil, err
	}

	if data != nil {
		if err := g.bindData(data, secret.GetSecret()); err != nil {
			return nil, err
		}
	}

	return secret, nil
}

func (g *GrpcPrivateHandler) bindData(src *structpb.Struct, dest any) error {
	config := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           dest,
		TagName:          "json",
		WeaklyTypedInput: true,
	}

	decoder, _ := mapstructure.NewDecoder(config)

	if err := decoder.Decode(src.AsMap()); err != nil {
		return status.Error(codes.InvalidArgument, "данные не соответствуют схеме")
	}

	return nil
}

func (g *GrpcPrivateHandler) mapToProto(results []any) (*structpb.ListValue, error) {
	payload, err := json.Marshal(results)

	if err != nil {
		return nil, status.Error(codes.Internal, "ошибка сериализации данных секрета")
	}

	var raw []interface{}

	if err := json.Unmarshal(payload, &raw); err != nil {
		return nil, status.Error(codes.Internal, "ошибка подготовки данных для proto")
	}

	list, err := structpb.NewList(raw)

	if err != nil {
		return nil, status.Error(codes.Internal, "ошибка формирования ответа (конвертация в proto)")
	}

	return list, nil
}

func (g *GrpcPrivateHandler) mergeAdditionalData(secretType string, data *structpb.Struct) {
	if data == nil || data.Fields == nil {
		return
	}

	if secretType == "Card" {
		if numberVal, ok := data.Fields["number"]; ok {
			cardNumber := numberVal.GetStringValue()
			cardType := secrets.GetCardType(cardNumber)
			data.Fields["card_type"] = structpb.NewStringValue(cardType)
			data.Fields["number"] = structpb.NewStringValue(secrets.FormatCardNumber(cardNumber))
		}
	}
}
