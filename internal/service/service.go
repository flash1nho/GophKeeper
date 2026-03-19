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
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/flash1nho/GophKeeper/internal/config"

	"go.uber.org/zap"
)

type Service struct {
	gHandler          *pb.GrpcHandler
	GrpcServerAddress string
	log               *zap.Logger
}

func NewService(gHandler *pb.GrpcHandler, settings config.SettingsObject) *Service {
	return &Service{
		gHandler:          gHandler,
		GrpcServerAddress: settings.GrpcServerAddress,
		log:               settings.Log,
	}
}

func runGrpcServer(ctx context.Context, s *Service) {
	serverErr := make(chan error, 1)
	creds := insecure.NewCredentials()

	grpcServer := grpc.NewServer(
		grpc.Creds(creds),
	)

	go func() {
		listen, err := net.Listen("tcp", s.GrpcServerAddress)

		if err == nil {
			pb.RegisterGophKeeperServiceServer(grpcServer, s.gHandler)

			s.log.Info("сервер gRPC начал работу")

			if err := grpcServer.Serve(listen); err != nil {
				s.log.Error("Ошибка при работе gRPC сервера", zap.Error(err))
			}
		} else {
			s.log.Error("Ошибка при инициализации gRPC listener", zap.Error(err))
		}
	}()

	select {
	case err := <-serverErr:
		s.log.Error(err.Error())
	case <-ctx.Done():
		s.log.Info("Завершение работы gRPC сервера")

		grpcServer.GracefulStop()

		s.log.Info("gRPC сервер успешно остановлен")
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
		s.log.Error("Работа завершена с ошибкой", zap.Error(err))
	}

	s.log.Info("Сохранение данных в хранилище...")

	s.gHandler.Pool.Close()

	s.log.Info("Все серверы успешно завершили работу.")
}
