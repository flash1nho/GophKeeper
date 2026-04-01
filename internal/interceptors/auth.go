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
	"github.com/flash1nho/GophKeeper/internal/security"
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

		userID, err := getUserIDFromToken(ctx, pool, settings, values[0])

		if err != nil {
			return nil, err
		}

		withValueCtx := context.WithValue(ctx, userKey, userID)

		return handler(withValueCtx, req)
	}
}

func StreamInterceptorAuth(pool *pgxpool.Pool, settings config.SettingsObject) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		md, ok := metadata.FromIncomingContext(ctx)

		if !ok {
			return status.Error(codes.Unauthenticated, "метаданные отсутствуют")
		}

		values := md["authorization"]

		if len(values) == 0 || values[0] == "" {
			return handler(srv, ss)
		}

		userID, err := getUserIDFromToken(ctx, pool, settings, values[0])

		if err != nil {
			return err
		}

		wrapped := &wrappedStream{
			ServerStream: ss,
			ctx:          context.WithValue(ctx, userKey, userID),
		}

		return handler(srv, wrapped)
	}
}

type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedStream) Context() context.Context {
	return w.ctx
}

func GetUserIDFromContext(ctx context.Context) (int, error) {
	userID, ok := ctx.Value(userKey).(int)

	if !ok {
		return 0, ErrUserNotFound
	}

	return userID, nil
}

func getUserIDFromToken(ctx context.Context, pool *pgxpool.Pool, settings config.SettingsObject, tokenParam string) (int, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenParam, jwt.MapClaims{})

	if err != nil {
		settings.Log.Error(err.Error())

		return 0, status.Error(codes.Unauthenticated, "недопустимый формат токена")
	}

	var userID int

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		manager := security.NewCryptoManager(settings.MasterKey)
		userID, err = manager.GetUserIDFromToken(claims)

		if err != nil {
			settings.Log.Error(err.Error())
		}
	}

	if userID == 0 {
		return 0, status.Error(codes.Unauthenticated, "недействительные требования")
	}

	err = users.UserVerify(ctx, pool, settings, userID, tokenParam)

	if err != nil {
		return 0, status.Error(codes.Unauthenticated, "Верификация не пройдена")
	}

	return userID, nil
}
