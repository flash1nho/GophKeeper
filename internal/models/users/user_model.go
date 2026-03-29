package users

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/flash1nho/GophKeeper/config"
	"github.com/flash1nho/GophKeeper/internal/security"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrCredentialsRequired = errors.New("введите логин и пароль")
	ErrSecretRequired      = errors.New("введите секретное слово")
	ErrUserAlreadyExists   = errors.New("логин уже существует")
	ErrUserNotFound        = errors.New("логин не существует")
	ErrInvalidPassword     = errors.New("неверный пароль")
	ErrInvalidToken        = errors.New("невалидный токен")
)

type User struct {
	ID       int
	Login    string
	Password string
	Secret   []byte
}

func NewUser(id int, login string, password string, secret string) *User {
	return &User{
		ID:       id,
		Login:    login,
		Password: password,
		Secret:   []byte(secret),
	}
}

func (user *User) UserRegister(ctx context.Context, pool *pgxpool.Pool, settings config.SettingsObject) (string, error) {
	if user.Login == "" || user.Password == "" {
		return "", ErrCredentialsRequired
	}

	if len(user.Secret) == 0 {
		return "", ErrSecretRequired
	}

	manager := security.NewCryptoManager(settings.MasterKey)
	passwordHash, err := manager.HashPassword(user.Password)

	if err != nil {
		return "", err
	}

	encryptedSecret, err := manager.Encrypt(user.Secret)

	if err != nil {
		return "", err
	}

	query, args, err := squirrel.Insert("users").
		Columns("login", "password_hash", "encrypted_secret", "created_at").
		Values(user.Login, passwordHash, encryptedSecret, time.Now().UTC()).
		Suffix("RETURNING id, encrypted_secret").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return "", err
	}

	err = pool.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Secret)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return "", ErrUserAlreadyExists
		}

		return "", err
	}

	return manager.GenerateToken(user.ID, user.Secret)
}

func (user *User) UserLogin(ctx context.Context, pool *pgxpool.Pool, settings config.SettingsObject) (string, error) {
	if user.Login == "" || user.Password == "" {
		return "", ErrCredentialsRequired
	}

	inputPassword := user.Password

	query, args, err := squirrel.Select("id", "password_hash", "encrypted_secret").
		From("users").
		Where(squirrel.Eq{"login": user.Login}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return "", err
	}

	err = pool.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Password, &user.Secret)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrUserNotFound
		}

		return "", err
	}

	manager := security.NewCryptoManager(settings.MasterKey)

	if !manager.CheckPassword(inputPassword, user.Password) {
		return "", ErrInvalidPassword
	}

	return manager.GenerateToken(user.ID, user.Secret)
}

func (user *User) UserVerify(ctx context.Context, pool *pgxpool.Pool, settings config.SettingsObject, token string) error {
	query, args, err := squirrel.Select("id", "encrypted_secret").
		From("users").
		Where(squirrel.Eq{"id": user.ID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	err = pool.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Secret)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}

		return err
	}

	manager := security.NewCryptoManager(settings.MasterKey)
	userID, err := manager.ValidateToken(token, user.Secret)

	if err != nil {
		return err
	}

	if userID != user.ID {
		return ErrInvalidToken
	}

	return nil
}
