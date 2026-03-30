package secrets

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"time"

	"dario.cat/mergo"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
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
	UpdatedAt     time.Time               `json:"-"`
	CryptoManager *security.CryptoManager `json:"-"`
}

type SecretResponse struct {
	ID        int    `json:"id"`
	Data      any    `json:"data"`
	Type      string `json:"type"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Secret interface {
	GetBaseSecret() *BaseSecret
	GetType() string
	GetSecret() any
	CreateValidate() error
	UpdateValidate() error
}

func (b *BaseSecret) GetBaseSecret() *BaseSecret {
	return b
}

func Create(ctx context.Context, pool *pgxpool.Pool, s Secret) ([]any, error) {
	err := s.CreateValidate()

	if err != nil {
		return nil, err
	}

	baseSecret := s.GetBaseSecret()

	if baseSecret.UserID == 0 {
		return nil, ErrInvalidUserID
	}

	payload, err := json.Marshal(s.GetSecret())

	if err != nil {
		return nil, err
	}

	userKey, err := baseSecret.getUserKey(ctx, pool)
	encryptedData, err := baseSecret.encryptData(payload, userKey)

	if err != nil {
		return nil, err
	}

	dateAt := time.Now().UTC()

	query, args, err := squirrel.Insert("secrets").
		Columns("user_id", "encrypted_data", "type", "created_at", "updated_at").
		Values(baseSecret.UserID, encryptedData, s.GetType(), dateAt, dateAt).
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	err = pool.QueryRow(ctx, query, args...).Scan(&baseSecret.ID)

	if err != nil {
		return nil, err
	}

	query, args, err = squirrel.Select("id", "encrypted_data", "created_at", "updated_at").
		From("secrets").
		Where(squirrel.Eq{"id": baseSecret.ID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	return baseSecret.data(ctx, pool, s, query, args)
}

func Get(ctx context.Context, pool *pgxpool.Pool, s Secret, ID int) ([]any, error) {
	if ID == 0 {
		return nil, ErrIDEmpty
	}

	baseSecret := s.GetBaseSecret()

	query, args, err := squirrel.Select("id", "encrypted_data", "created_at", "updated_at").
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

	query, args, err := squirrel.Select("id", "encrypted_data", "created_at", "updated_at").
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

func Update(ctx context.Context, pool *pgxpool.Pool, s Secret, ID int) ([]any, error) {
	if ID == 0 {
		return nil, ErrIDEmpty
	}

	err := s.UpdateValidate()

	if err != nil {
		return nil, err
	}

	baseSecret := s.GetBaseSecret()
	userKey, err := baseSecret.getUserKey(ctx, pool)

	tx, err := pool.Begin(ctx)

	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	sqlSelect, args, _ := squirrel.Select("encrypted_data").
		From("secrets").
		Where(squirrel.Eq{"id": ID}).
		Where(squirrel.Eq{"user_id": baseSecret.UserID}).
		Where(squirrel.Eq{"type": s.GetType()}).
		Suffix("FOR UPDATE SKIP LOCKED").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	var encryptedData []byte

	if err := tx.QueryRow(ctx, sqlSelect, args...).Scan(&encryptedData); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrEmptyRows
		}

		return nil, err
	}

	if err != nil {
		return nil, err
	}

	oldPayload, err := baseSecret.decryptData(encryptedData, userKey)

	if err != nil {
		return nil, err
	}

	var oldMap map[string]any

	if err := json.Unmarshal(oldPayload, &oldMap); err != nil {
		return nil, err
	}

	newData := s.GetSecret()
	newPayload, err := json.Marshal(newData)

	if err != nil {
		return nil, err
	}

	var newMap map[string]any

	if err := json.Unmarshal(newPayload, &newMap); err != nil {
		return nil, err
	}

	for k, v := range newMap {
		if v == nil || v == "" || v == 0 {
			delete(newMap, k)
		}
	}

	if err := mergo.Merge(&oldMap, newMap, mergo.WithOverride); err != nil {
		return nil, err
	}

	payload, err := json.Marshal(oldMap)

	if err != nil {
		return nil, err
	}

	newEncryptedData, err := baseSecret.encryptData(payload, userKey)

	if err != nil {
		return nil, err
	}

	sqlUpdate, args, err := squirrel.Update("secrets").
		Set("encrypted_data", newEncryptedData).
		Set("updated_at", time.Now().UTC()).
		Where(squirrel.Eq{"id": ID}).
		Where(squirrel.Eq{"user_id": baseSecret.UserID}).
		Where(squirrel.Eq{"type": s.GetType()}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(ctx, sqlUpdate, args...)

	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)

	if err != nil {
		return nil, err
	}

	query, args, err := squirrel.Select("id", "encrypted_data", "created_at", "updated_at").
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

func Delete(ctx context.Context, pool *pgxpool.Pool, s Secret, ID int) error {
	if ID == 0 {
		return ErrIDEmpty
	}

	baseSecret := s.GetBaseSecret()

	query, args, err := squirrel.Delete("secrets").
		Where(squirrel.Eq{"id": ID}).
		Where(squirrel.Eq{"user_id": baseSecret.UserID}).
		Where(squirrel.Eq{"type": s.GetType()}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	_, err = pool.Exec(ctx, query, args...)

	if err != nil {
		return err
	}

	return nil
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
		var updatedAt time.Time

		if err := rows.Scan(&id, &encryptedData, &createdAt, &updatedAt); err != nil {
			return nil, err
		}

		decryptedData, err := baseSecret.decryptData(encryptedData, userKey)

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
			Type:      s.GetType(),
			CreatedAt: createdAt.Format("02.01.2006 15:04:05"),
			UpdatedAt: updatedAt.Format("02.01.2006 15:04:05"),
		}

		results = append(results, result)
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return baseSecret.CryptoManager.Decrypt(encryptedSecret, baseSecret.CryptoManager.MasterKey)
}

func (baseSecret *BaseSecret) decryptData(encryptedData []byte, userKey []byte) ([]byte, error) {
	decryptedUserData, err := baseSecret.CryptoManager.Decrypt(encryptedData, baseSecret.CryptoManager.MasterKey)

	if err != nil {
		return nil, err
	}

	return baseSecret.CryptoManager.Decrypt(decryptedUserData, userKey)
}

func (baseSecret *BaseSecret) encryptData(data []byte, userKey []byte) ([]byte, error) {
	encryptedUserData, err := baseSecret.CryptoManager.Encrypt(data, userKey)

	if err != nil {
		return nil, err
	}

	return baseSecret.CryptoManager.Encrypt(encryptedUserData, baseSecret.CryptoManager.MasterKey)
}
