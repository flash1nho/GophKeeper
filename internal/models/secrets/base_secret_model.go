package secrets

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BaseSecret struct {
	ID        uint32
	UserID    int
	CreatedAt time.Time
}

type TextNote struct {
	BaseSecret

	Content string
}

type Secret interface {
	GetBaseSecret() *BaseSecret
	GetType() string
	GetPayload() any
	Validate() error
}

func Create(ctx context.Context, pool *pgxpool.Pool, s Secret) error {
	err := s.Validate()

	if err != nil {
		return err
	}

	baseSecret := s.GetBaseSecret()

	if baseSecret.UserID == 0 {
		return fmt.Errorf("UserID не может быть пустынм")
	}

	payload, err := json.Marshal(s.GetPayload())

	if err != nil {
		return err
	}

	query, args, err := squirrel.Insert("secrets").
		Columns("user_id", "properties", "type", "created_at").
		Values(baseSecret.UserID, payload, s.GetType(), time.Now().UTC()).
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	return pool.QueryRow(ctx, query, args...).Scan(&baseSecret.ID)
}
