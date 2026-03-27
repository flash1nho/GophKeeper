package grpc

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/flash1nho/GophKeeper/internal/models/secrets"
)

type GrpcPrivateHandler struct {
	UnimplementedGophKeeperPrivateServiceServer

	Pool *pgxpool.Pool
	Log  *zap.Logger
}

func (g *GrpcPrivateHandler) TextNoteCreate(ctx context.Context, req *TextNoteCreateRequest) (*CreateResponse, error) {
	var response CreateResponse

	textNote := secrets.NewTextNote(1, req.Content)
	err := secrets.Create(ctx, g.Pool, textNote)

	if err != nil {
		return nil, err
	}

	response.ID = textNote.ID

	return &response, nil
}
