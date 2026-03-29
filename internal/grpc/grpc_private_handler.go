package grpc

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/flash1nho/GophKeeper/config"
	"github.com/flash1nho/GophKeeper/internal/facade"
	"github.com/flash1nho/GophKeeper/internal/models/secrets"
)

type GrpcPrivateHandler struct {
	UnimplementedGophKeeperPrivateServiceServer

	Pool     *pgxpool.Pool
	Settings config.SettingsObject
	facade   *facade.Facade
}

func (g *GrpcPrivateHandler) TextCreate(ctx context.Context, req *TextCreateRequest) (*CreateResponse, error) {
	var response CreateResponse

	userID, err := g.facade.GetUserFromContext(ctx)

	if err != nil {
		return nil, err
	}

	textNote := secrets.NewText(userID, req.Content)
	err = secrets.Create(ctx, g.Pool, textNote)

	if err != nil {
		return nil, err
	}

	response.ID = int32(textNote.ID)

	return &response, nil
}
