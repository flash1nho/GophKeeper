package secrets

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"reflect"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/flash1nho/GophKeeper/internal/security"
)

var (
	ErrInvalidUserID  = errors.New("id пользователя не может быть пустым")
	ErrUserNotFound   = errors.New("пользователь не существует")
	ErrSecretNotFound = errors.New("пользовательские данные не существуют")
	ErrUnknownType    = errors.New("тип не найден")
	ErrEmptyRows      = errors.New("секреты не найдены")
	ErrIDEmpty        = errors.New("'id' не может быть пустым")
)

type BaseSecret struct {
	ID            int                     `json:"-"`
	UserID        int                     `json:"-"`
	CreatedAt     time.Time               `json:"-"`
	CryptoManager *security.CryptoManager `json:"-"`
}

type SecretResponse struct {
	ID        int    `json:"id"`
	Data      any    `json:"data"`
	CreatedAt string `json:"created_at"`
}

type Secret interface {
	GetBaseSecret() *BaseSecret
	GetType() string
	GetSecret() any
	CreateValidate() error
	UpdateValidate() error
}

func Create(ctx context.Context, pool *pgxpool.Pool, s Secret) error {
	err := s.CreateValidate()

	if err != nil {
		return err
	}

	baseSecret := s.GetBaseSecret()

	if baseSecret.UserID == 0 {
		return ErrInvalidUserID
	}

	payload, err := json.Marshal(s.GetSecret())

	if err != nil {
		return err
	}

	userKey, err := baseSecret.getUserKey(ctx, pool)

	encryptedUserData, err := baseSecret.CryptoManager.Encrypt(payload, userKey)

	if err != nil {
		return err
	}

	encryptedData, err := baseSecret.CryptoManager.Encrypt(encryptedUserData, baseSecret.CryptoManager.MasterKey)

	if err != nil {
		return err
	}

	query, args, err := squirrel.Insert("secrets").
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

func Get(ctx context.Context, pool *pgxpool.Pool, s Secret, ID int) ([]any, error) {
	baseSecret := s.GetBaseSecret()

	if ID == 0 {
		return nil, ErrIDEmpty
	}

	query, args, err := squirrel.Select("id", "encrypted_data", "created_at").
		From("secrets").
		Where(squirrel.Eq{"id": ID}).
		Where(squirrel.Eq{"user_id": baseSecret.UserID}).
		Where(squirrel.Eq{"type": s.GetType()}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	return baseSecret.data(ctx, pool, s, query, args)
}

func List(ctx context.Context, pool *pgxpool.Pool, s Secret) ([]any, error) {
	baseSecret := s.GetBaseSecret()

	query, args, err := squirrel.Select("id", "encrypted_data", "created_at").
		From("secrets").
		Where(squirrel.Eq{"user_id": baseSecret.UserID}).
		Where(squirrel.Eq{"type": s.GetType()}).
		OrderBy("id ASC").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	return baseSecret.data(ctx, pool, s, query, args)
}

func (baseSecret *BaseSecret) data(ctx context.Context, pool *pgxpool.Pool, s Secret, query string, args []interface{}) ([]any, error) {
	userKey, err := baseSecret.getUserKey(ctx, pool)

	if err != nil {
		return nil, err
	}

	rows, err := pool.Query(ctx, query, args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []any

	secretObject := s.GetSecret()

	if secretObject == nil {
		return nil, ErrUnknownType
	}

	targetType := reflect.TypeOf(secretObject)

	if targetType.Kind() == reflect.Ptr {
		targetType = targetType.Elem()
	}

	for rows.Next() {
		var id int
		var encryptedData []byte
		var createdAt time.Time

		if err := rows.Scan(&id, &encryptedData, &createdAt); err != nil {
			return nil, err
		}

		encryptedUserData, err := baseSecret.CryptoManager.Decrypt(encryptedData, baseSecret.CryptoManager.MasterKey)

		if err != nil {
			return nil, err
		}

		decryptedData, err := baseSecret.CryptoManager.Decrypt(encryptedUserData, userKey)

		if err != nil {
			return nil, err
		}

		secret := reflect.New(targetType).Interface()

		if err := json.Unmarshal(decryptedData, secret); err != nil {
			return nil, err
		}

		result := SecretResponse{
			ID:        id,
			Data:      secret,
			CreatedAt: createdAt.Format("02.01.2006 15:04"),
		}

		results = append(results, result)
	}

	if len(results) == 0 {
		return nil, ErrEmptyRows
	}

	return results, rows.Err()
}

func Update(ctx context.Context, pool *pgxpool.Pool, s Secret) error {
	err := s.UpdateValidate()

	if err != nil {
		return err
	}

	payload, err := json.Marshal(s.GetSecret())

	if err != nil {
		return err
	}

	baseSecret := s.GetBaseSecret()
	userKey, err := baseSecret.getUserKey(ctx, pool)

	encryptedUserData, err := baseSecret.CryptoManager.Encrypt(payload, userKey)

	if err != nil {
		return err
	}

	encryptedData, err := baseSecret.CryptoManager.Encrypt(encryptedUserData, baseSecret.CryptoManager.MasterKey)

	if err != nil {
		return err
	}

	query, args, err := squirrel.Update("secrets").
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

func (baseSecret *BaseSecret) getUserKey(ctx context.Context, pool *pgxpool.Pool) ([]byte, error) {
	var encryptedSecret []byte

	query, args, err := squirrel.Select("encrypted_secret").
		From("users").
		Where(squirrel.Eq{"id": baseSecret.UserID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	err = pool.QueryRow(ctx, query, args...).Scan(&encryptedSecret)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return baseSecret.CryptoManager.Decrypt(encryptedSecret, baseSecret.CryptoManager.MasterKey)
}
