package users

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int
	Login    string
	Password string
	Secret   string
}

func NewUser(login string, password string) *User {
	return &User{
		Login:    login,
		Password: password,
	}
}

func (user *User) UserRegister(ctx context.Context, pool *pgxpool.Pool) (string, error) {
	if user.Login == "" || user.Password == "" {
		return "", fmt.Errorf("введите логин и пароль")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	query, args, err := squirrel.Insert("users").
		Columns("login", "password", "created_at").
		Values(user.Login, hashedPassword, time.Now().UTC()).
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return "", err
	}

	err = pool.QueryRow(ctx, query, args...).Scan(&user.ID)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return "", fmt.Errorf("логин уже существует")
		}

		return "", err
	}

	return user.createToken()
}

func (user *User) UserLogin(ctx context.Context, pool *pgxpool.Pool) (string, error) {
	if user.Login == "" || user.Password == "" {
		return "", fmt.Errorf("введите логин и пароль")
	}

	password := user.Password

	query, args, err := squirrel.Select("id", "password").
		From("users").
		Where(squirrel.Eq{"login": user.Login}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return "", err
	}

	err = pool.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Password)

	if err != nil {
		return "", fmt.Errorf("Пользователь не найден")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", fmt.Errorf("Неверный пароль")
	}

	return user.createToken()
}

func (user *User) createToken() (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(user.Secret))
}
