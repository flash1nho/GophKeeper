package cmd

import (
	pb "github.com/flash1nho/GophKeeper/internal/grpc"
	"google.golang.org/grpc"
)

func PrivateClient() (pb.GophKeeperPrivateServiceClient, *grpc.ClientConn, error) {
	conn, err := BaseConnection()

	if err != nil {
		return nil, nil, err
	}

	client := pb.NewGophKeeperPrivateServiceClient(conn)

	return client, conn, nil
}
