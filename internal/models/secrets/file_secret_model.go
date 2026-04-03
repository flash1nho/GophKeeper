package secrets

import (
	"context"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/flash1nho/GophKeeper/internal/security"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"errors"
)

var (
	ErrEmptyFileName = errors.New("название файла не может быть пустым")
)

type File struct {
	BaseSecret

	FileName        string `json:"file_name"`
	FileContentType string `json:"file_content_type"`
}

func NewFile(userID int, masterKey []byte, pool *pgxpool.Pool) *File {
	return &File{
		BaseSecret: BaseSecret{
			UserID:        userID,
			CryptoManager: security.NewCryptoManager(masterKey),
			pool:          pool,
		},
	}
}

func (s *File) GetType() string {
	return "File"
}

func (s *File) GetSecret() any {
	return s
}

func (s *File) CreateValidate(ctx context.Context) error {
	if strings.TrimSpace(s.FileName) == "" {
		return ErrEmptyFileName
	}

	return nil
}

func (s *File) UpdateValidate() error {
	return nil
}

func (s *File) FileExists(ctx context.Context) (bool, error) {
	baseSecret := s.GetBaseSecret()

	builder := squirrel.Select("id", "file_offset").
		From("secrets").
		Where(squirrel.Eq{"user_id": baseSecret.UserID}).
		Where(squirrel.Eq{"type": s.GetType()})

	if baseSecret.ID > 0 {
		builder = builder.Where(squirrel.Eq{"id": baseSecret.ID})
	} else {
		builder = builder.Where(squirrel.Eq{"file_name": baseSecret.FileName})
	}

	query, args, err := builder.PlaceholderFormat(squirrel.Dollar).ToSql()

	if err != nil {
		return false, err
	}

	err = baseSecret.pool.QueryRow(ctx, query, args...).Scan(&baseSecret.ID, &baseSecret.FileOffset)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	if baseSecret.ID == 0 {
		return false, nil
	}

	return true, nil
}
