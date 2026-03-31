package service

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	pb "github.com/flash1nho/GophKeeper/internal/grpc"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"

	"github.com/flash1nho/GophKeeper/internal/interceptors"
	"google.golang.org/grpc/credentials"

	"github.com/flash1nho/GophKeeper/certs"

	"go.uber.org/zap"
)

type Service struct {
	Handler *pb.GrpcHandler
}

func NewService(handler *pb.GrpcHandler) *Service {
	return &Service{
		Handler: handler,
	}
}

func runGrpcServer(ctx context.Context, s *Service) {
	serverErr := make(chan error, 1)
	certs, err := certs.NewCerts("server")

	if err != nil {
		s.Handler.Settings.Log.Error("Сертификаты не найдены", zap.Error(err))
	}

	creds, err := credentials.NewServerTLSFromFile(certs.Cert, certs.Key)

	if err != nil {
		s.Handler.Settings.Log.Error("Не удалось настроить TLS", zap.Error(err))
	}

	loggingInterceptor := logging.UnaryServerInterceptor(interceptors.InterceptorLogger(s.Handler.Settings.Log))
	authInterceptor := interceptors.InterceptorAuth(s.Handler.Pool, s.Handler.Settings)
	authStreamInterceptor := interceptors.StreamInterceptorAuth(s.Handler.Pool, s.Handler.Settings)

	grpcServer := grpc.NewServer(
		grpc.Creds(creds),
		grpc.ChainUnaryInterceptor(
			loggingInterceptor,
			authInterceptor,
		),
		grpc.StreamInterceptor(authStreamInterceptor),
	)

	go func() {
		listen, err := net.Listen("tcp", s.Handler.Settings.GrpcServerAddress)

		if err == nil {
			pb.RegisterGophKeeperPublicServiceServer(grpcServer, s.Handler.GrpcPublicHandler)
			pb.RegisterGophKeeperPrivateServiceServer(grpcServer, s.Handler.GrpcPrivateHandler)

			s.Handler.Settings.Log.Info("сервер gRPC начал работу")

			if err := grpcServer.Serve(listen); err != nil {
				s.Handler.Settings.Log.Error("Ошибка при работе gRPC сервера", zap.Error(err))
			}
		} else {
			s.Handler.Settings.Log.Error("Ошибка при инициализации gRPC listener", zap.Error(err))
		}
	}()

	select {
	case err := <-serverErr:
		s.Handler.Settings.Log.Error(err.Error())
	case <-ctx.Done():
		s.Handler.Settings.Log.Info("Завершение работы gRPC сервера")

		grpcServer.GracefulStop()

		s.Handler.Settings.Log.Info("gRPC сервер успешно остановлен")
	}
}

func (s *Service) Run() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		runGrpcServer(ctx, s)
		return nil
	})

	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		s.Handler.Settings.Log.Error("Работа завершена с ошибкой", zap.Error(err))
	}

	s.Handler.Settings.Log.Info("Сохранение данных в хранилище...")

	s.Handler.Pool.Close()

	s.Handler.Settings.Log.Info("Все серверы успешно завершили работу.")
}
