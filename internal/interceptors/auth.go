package interceptors

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/flash1nho/GophKeeper/config"
	"github.com/flash1nho/GophKeeper/internal/models/users"
)

var (
	ErrUserNotFound = errors.New("неверный токен")
)

type UserID string

const userKey = UserID("userID")

func InterceptorAuth(pool *pgxpool.Pool, settings config.SettingsObject) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)

		if !ok {
			return nil, status.Error(codes.Unauthenticated, "метаданные отсутствуют")
		}

		values := md["authorization"]

		if len(values) == 0 || values[0] == "" {
			return handler(ctx, req)
		}

		tokenParam := values[0]
		token, _, err := new(jwt.Parser).ParseUnverified(tokenParam, jwt.MapClaims{})

		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "недопустимый формат токена")
		}

		user := users.NewUser(0, "", "", "")

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if sub, ok := claims["sub"].(float64); ok {
				user.ID = int(sub)
			}
		}

		if user.ID == 0 {
			return nil, status.Error(codes.Unauthenticated, "недействительные требования")
		}

		err = user.UserVerify(ctx, pool, settings, tokenParam)

		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "Верификация не пройдена")
		}

		withValueCtx := context.WithValue(ctx, userKey, user.ID)

		return handler(withValueCtx, req)
	}
}

func GetUserFromContext(ctx context.Context) (int, error) {
	userID, ok := ctx.Value(userKey).(int)

	if !ok {
		return 0, ErrUserNotFound
	}

	return userID, nil
}
