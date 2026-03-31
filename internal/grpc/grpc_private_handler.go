package grpc

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/flash1nho/GophKeeper/config"
	"github.com/flash1nho/GophKeeper/internal/facade"
	"github.com/flash1nho/GophKeeper/internal/models/secrets"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	FileStoragePath = "./uploads"
	ChunkSize       = 1024 * 1024 // 1MB
)

var typeToSecret = map[string]func(userID int, masterKey []byte, pool *pgxpool.Pool) secrets.Secret{
	"Text": func(u int, k []byte, p *pgxpool.Pool) secrets.Secret { return secrets.NewText(u, k, p) },
	"Cred": func(u int, k []byte, p *pgxpool.Pool) secrets.Secret { return secrets.NewCred(u, k, p) },
	"Card": func(u int, k []byte, p *pgxpool.Pool) secrets.Secret { return secrets.NewCard(u, k, p) },
	"File": func(u int, k []byte, p *pgxpool.Pool) secrets.Secret { return secrets.NewFile(u, k, p) },
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

	res, err := secrets.Create(ctx, secret)

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

	res, err := secrets.Get(ctx, secret, int(req.ID))

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

	res, err := secrets.List(ctx, secret)

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

	res, err := secrets.Update(ctx, secret, int(req.ID))

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

	if err := secrets.Delete(ctx, secret, int(req.ID)); err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка при удалении секрета: %v", err)
	}

	return &DeleteResponse{}, nil
}

func (g *GrpcPrivateHandler) Upload(stream GophKeeperPrivateService_UploadServer) error {
	ctx := stream.Context()
	secret, err := g.getSecretInstance(ctx, "File")

	if err != nil {
		return err
	}

	baseSecret := secret.GetBaseSecret()
	userKey, err := baseSecret.GetUserKey(ctx)

	if err != nil {
		return err
	}

	var file *os.File
	var currentOffset int64

	for {
		req, err := stream.Recv()

		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		if meta := req.GetMetadata(); meta != nil {
			baseSecret.FileName = meta.FileName
			fileExists, err := secret.FileExists(ctx)

			// TODO

			if fileExists {
				bits, err := json.Marshal(meta)

				if err != nil {
					return err
				}

				metaStruct := &structpb.Struct{}

				if err := metaStruct.UnmarshalJSON(bits); err != nil {
					return err
				}

				secret, err = g.prepareSecret(ctx, "File", metaStruct)

				if err != nil {
					return err
				}

				if _, err := secrets.Create(ctx, secret); err != nil {
					return status.Errorf(codes.Internal, "ошибка создания записи: %v", err)
				}
			}

			baseSecret.ID = secret.GetBaseSecret().ID

			if baseSecret.ID == 0 {
				return status.Error(codes.Internal, "ID секрета не определен")
			}

			filePath := filepath.Join(FileStoragePath, string(baseSecret.ID))
			file, err = os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)

			if err != nil {
				return err
			}

			defer file.Close()

			info, err := file.Stat()

			if err != nil {
				return err
			}

			currentOffset = info.Size()

			if meta.FileOffset != currentOffset {
				return status.Errorf(codes.Aborted, "разрыв данных: на сервере %d, клиент прислал %d", currentOffset, meta.FileOffset)
			}

			if _, err := file.Seek(currentOffset, 0); err != nil {
				return err
			}

			continue
		}

		if file == nil {
			return status.Error(codes.FailedPrecondition, "метаданные не получены")
		}

		chunk := req.GetChunk()

		if chunk == nil {
			continue
		}

		encryptedChunk, err := baseSecret.EncryptStream(chunk, userKey, currentOffset)

		if err != nil {
			return err
		}

		n, err := file.Write(encryptedChunk)

		if err != nil {
			return err
		}

		currentOffset += int64(n)
	}

	data := map[string]interface{}{"FileOffset": currentOffset}
	structData, err := structpb.NewStruct(data)

	if err != nil {
		return status.Errorf(codes.Internal, "ошибка конвертации метаданных: %v", err)
	}

	secret, err = g.prepareSecret(ctx, "File", structData)

	if err != nil {
		return err
	}

	res, err := secrets.Get(ctx, secret, baseSecret.ID)

	if err != nil {
		return status.Errorf(codes.Internal, "ошибка при обновлении секрета: %v", err)
	}

	protoData, err := g.mapToProto(res)

	if err != nil {
		return err
	}

	return stream.SendAndClose(&UploadResponse{Secrets: protoData})
}

// func (s *server) Download(req *DownloadRequest, stream FileService_DownloadServer) error {
// 	if err := authInterceptor(stream.Context()); err != nil {
// 		return err
// 	}

// 	// 1. Получить путь из БД по file_id
// 	filePath := filepath.Join(FileStoragePath, req.FileID)
// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	// 2. Применить докачку
// 	_, err = file.Seek(req.Offset, 0)
// 	if err != nil {
// 		return err
// 	}

// 	buffer := make([]byte, ChunkSize)
// 	for {
// 		bytesRead, err := file.Read(buffer)
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			return err
// 		}
// 		if err := stream.Send(&DownloadResponse{Chunk: buffer[:bytesRead]}); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

func (g *GrpcPrivateHandler) getSecretInstance(ctx context.Context, secretType string) (secrets.Secret, error) {
	userID, err := g.facade.GetUserIDFromContext(ctx)

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "пользователь не авторизован")
	}

	fn, ok := typeToSecret[secretType]

	if !ok {
		return nil, status.Error(codes.InvalidArgument, "неизвестный тип секрета")
	}

	return fn(userID, g.settings.MasterKey, g.pool), nil
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
