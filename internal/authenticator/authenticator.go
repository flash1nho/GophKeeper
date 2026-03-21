package authenticator

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type UserID string

const userKey = UserID("userID")

type Authenticator struct{}

func NewAuthenticator() *Authenticator {
	return &Authenticator{}
}

type AuthProvider interface {
	SetToken(ctx context.Context, UserID int) error
	ParseToken(ctx context.Context, token string) error
}

func FromContext(ctx context.Context) (int, error) {
	userID, ok := ctx.Value(userKey).(int)

	if !ok {
		return 0, fmt.Errorf("userID не найден в контексте")
	}

	return userID, nil
}

func GetUserKey() UserID {
	return userKey
}

func (a *Authenticator) Authenticate(ctx context.Context, provider AuthProvider) (context.Context, error) {
	// TODO
	// value, err := provider.parseToken(ctx)

	// if err != nil {
	// 	return nil, err
	// }

	// userID, err := a.(value)

	// if err != nil {
	// 	return nil, err
	// }

	// p.SetToken(ctx, value)

	// return context.WithValue(ctx, userKey, userID), nil
	return ctx, nil
}

func (a *Authenticator) CreateToken(UserID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": UserID,
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte("your-secret-key"))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
