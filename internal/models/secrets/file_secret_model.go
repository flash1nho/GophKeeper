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
	ErrFileExists    = errors.New("файл существует")
	ErrEmptyFileName = errors.New("название файла не может быть пустым")
)

type File struct {
	BaseSecret

	FileContentType string `json:"file_content_type"`
	FileSize        string `json:"file_size"`
	FileOffset      string `json:"file_offset`
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
	fileName := s.GetBaseSecret().FileName

	if strings.TrimSpace(fileName) == "" {
		return ErrEmptyFileName
	}

	fileExists, err := s.FileExists(ctx)

	if err != nil {
		return err
	}

	if fileExists {
		return ErrFileExists
	}

	return nil
}

func (s *File) UpdateValidate() error {
	return ErrNotImplemented
}

func (s *File) FileExists(ctx context.Context) (bool, error) {
	baseSecret := s.GetBaseSecret()

	query, args, err := squirrel.Select("id").
		From("secrets").
		Where(squirrel.Eq{"file_name": baseSecret.FileName}).
		Where(squirrel.Eq{"user_id": baseSecret.UserID}).
		Where(squirrel.Eq{"type": s.GetType()}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return false, err
	}

	err = baseSecret.pool.QueryRow(ctx, query, args...).Scan(&baseSecret.ID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return true, nil
		}

		return false, err
	}

	if baseSecret.ID > 0 {
		return true, nil
	}

	return false, err
}
