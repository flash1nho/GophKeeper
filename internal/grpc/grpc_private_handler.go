package grpc

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/flash1nho/GophKeeper/config"
	"github.com/flash1nho/GophKeeper/internal/facade"
	"github.com/flash1nho/GophKeeper/internal/models/secrets"
	"google.golang.org/protobuf/types/known/structpb"
)

var (
	ErrUnknownType   = errors.New("тип не найден")
	ErrFailedToParse = errors.New("не удалось проанализировать секретные данные")
)

type GrpcPrivateHandler struct {
	UnimplementedGophKeeperPrivateServiceServer

	Pool     *pgxpool.Pool
	Settings config.SettingsObject
	facade   *facade.Facade
}

func (g *GrpcPrivateHandler) Create(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
	jsonData, err := json.Marshal(req.Data.AsMap())

	if err != nil {
		return nil, err
	}

	secretObject, err := g.getSecretObject(ctx, req.Type)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(jsonData, secretObject.GetSecret()); err != nil {
		return nil, err
	}

	results, err := secrets.Create(ctx, g.Pool, secretObject)

	if err != nil {
		return nil, err
	}

	protoSecrets, err := g.response(ctx, results)

	if err != nil {
		return nil, err
	}

	return &CreateResponse{Secrets: protoSecrets}, nil
}

func (g *GrpcPrivateHandler) Get(ctx context.Context, req *GetRequest) (*GetResponse, error) {
	secretObject, err := g.getSecretObject(ctx, req.Type)

	if err != nil {
		return nil, err
	}

	results, err := secrets.Get(ctx, g.Pool, secretObject, int(req.ID))

	if err != nil {
		return nil, err
	}

	protoSecrets, err := g.response(ctx, results)

	if err != nil {
		return nil, err
	}

	return &GetResponse{Secrets: protoSecrets}, nil
}

func (g *GrpcPrivateHandler) List(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	secretObject, err := g.getSecretObject(ctx, req.Type)

	if err != nil {
		return nil, err
	}

	results, err := secrets.List(ctx, g.Pool, secretObject)

	if err != nil {
		return nil, err
	}

	protoSecrets, err := g.response(ctx, results)

	if err != nil {
		return nil, err
	}

	return &ListResponse{Secrets: protoSecrets}, nil
}

func (g *GrpcPrivateHandler) Update(ctx context.Context, req *UpdateRequest) (*UpdateResponse, error) {
	jsonData, err := json.Marshal(req.Data.AsMap())

	if err != nil {
		return nil, err
	}

	secretObject, err := g.getSecretObject(ctx, req.Type)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(jsonData, secretObject.GetSecret()); err != nil {
		return nil, err
	}

	results, err := secrets.Update(ctx, g.Pool, secretObject, int(req.ID))

	if err != nil {
		return nil, err
	}

	protoSecrets, err := g.response(ctx, results)

	if err != nil {
		return nil, err
	}

	return &UpdateResponse{Secrets: protoSecrets}, nil
}

func (g *GrpcPrivateHandler) Delete(ctx context.Context, req *DeleteRequest) (*DeleteResponse, error) {
	secretObject, err := g.getSecretObject(ctx, req.Type)

	if err != nil {
		return nil, err
	}

	err = secrets.Delete(ctx, g.Pool, secretObject, int(req.ID))

	if err != nil {
		return nil, err
	}

	return &DeleteResponse{}, nil
}

func (g *GrpcPrivateHandler) getSecretObject(ctx context.Context, secretType string) (secrets.Secret, error) {
	userID, err := g.facade.GetUserIDFromContext(ctx)

	if err != nil {
		return nil, err
	}

	var secretObject secrets.Secret

	switch secretType {
	case "Text":
		secretObject = secrets.NewText(userID, g.Settings)
	case "Cred":
		secretObject = secrets.NewCred(userID, g.Settings)
	default:
		return nil, ErrUnknownType
	}

	return secretObject, nil
}

func (g *GrpcPrivateHandler) response(ctx context.Context, results []any) (*structpb.ListValue, error) {
	payload, err := json.Marshal(results)

	if err != nil {
		return nil, err
	}

	var data []interface{}

	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, err
	}

	return structpb.NewList(data)
}
