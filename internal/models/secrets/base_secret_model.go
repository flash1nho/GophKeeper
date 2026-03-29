package secrets

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/flash1nho/GophKeeper/internal/models/users"
	"github.com/flash1nho/GophKeeper/internal/security"
)

var (
	ErrInvalidUserID = errors.New("id пользователя не может быть пустым")
	ErrUserNotFound  = errors.New("пользователь не существует")
)

type BaseSecret struct {
	ID        int
	UserID    int
	CreatedAt time.Time
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
		return ErrInvalidUserID
	}

	query, args, err := squirrel.Select("encrypted_secret").
		From("users").
		Where(squirrel.Eq{"id": baseSecret.UserID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	var user users.User

	err = pool.QueryRow(ctx, query, args...).Scan(&user.Secret)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}

		return err
	}

	payload, err := json.Marshal(s.GetPayload())

	if err != nil {
		return err
	}

	manager := security.NewCryptoManager(user.Secret)
	encryptedData, err := manager.Encrypt(payload)

	if err != nil {
		return err
	}

	query, args, err = squirrel.Insert("secrets").
		Columns("user_id", "encrypted_data", "type", "created_at").
		Values(baseSecret.UserID, encryptedData, s.GetType(), time.Now().UTC()).
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	return pool.QueryRow(ctx, query, args...).Scan(&baseSecret.ID)
}
