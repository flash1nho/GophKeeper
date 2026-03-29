package secrets

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
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
)

type BaseSecret struct {
	ID            int                     `json:"-"`
	UserID        int                     `json:"-"`
	CreatedAt     time.Time               `json:"-"`
	CryptoManager *security.CryptoManager `json:"-"`
}

type Secret interface {
	GetBaseSecret() *BaseSecret
	GetType() string
	GetSecret() any
	Validate() error
}

func Create(ctx context.Context, pool *pgxpool.Pool, s Secret) error {
	err := s.Validate()

	if err != nil {
		return err
	}

	fmt.Println(s.GetSecret())

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

	query, args, err := squirrel.Select("encrypted_data").
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

	query, args, err := squirrel.Select("encrypted_data").
		From("secrets").
		Where(squirrel.Eq{"user_id": baseSecret.UserID}).
		Where(squirrel.Eq{"type": s.GetType()}).
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
		var encryptedData []byte

		if err := rows.Scan(&encryptedData); err != nil {
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

		results = append(results, secret)
	}

	if len(results) == 0 {
		return nil, ErrEmptyRows
	}

	return results, rows.Err()
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
