package interceptors

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"

	"go.uber.org/zap"
)

func InterceptorLogger(l *zap.Logger) logging.Logger {
	return logging.LoggerFunc(func(_ context.Context, lvl logging.Level, msg string, fields ...any) {

		f := make([]zap.Field, 0, len(fields)/2)
		for i := 0; i < len(fields); i += 2 {
			key := fields[i].(string)
			value := fields[i+1]
			f = append(f, zap.Any(key, value))
		}

		l.Log(zap.InfoLevel, msg, f...)
	})
}
