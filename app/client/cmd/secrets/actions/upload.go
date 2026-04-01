package actions

import (
	"context"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/flash1nho/GophKeeper/app/client/helpers"
	"github.com/flash1nho/GophKeeper/config"

	pb "github.com/flash1nho/GophKeeper/internal/grpc"
)

func SecretsUploadCommand(client *pb.GophKeeperPrivateServiceClient, settings config.SettingsObject) *cobra.Command {
	cmd := &cobra.Command{
		Use: "upload",
		Run: func(cmd *cobra.Command, args []string) {
			filePath := args[0]
			file, err := os.Open(filePath)

			if err != nil {
				settings.Log.Fatal(err.Error())
			}

			stat, err := file.Stat()

			if err != nil {
				settings.Log.Fatal(err.Error())
			}

			statusResp, err := (*client).GetUploadStatus(context.Background(), &pb.UploadStatusRequest{
				FileName: stat.Name(),
			})

			if err != nil {
				settings.Log.Fatal(err.Error())
			}

			var remoteOffset int64 = 0

			if statusResp != nil {
				remoteOffset = statusResp.FileOffset
			}

			if remoteOffset >= stat.Size() {
				fmt.Println("✅ Файл уже загружен")
				return
			}

			_, err = file.Seek(remoteOffset, io.SeekStart)

			if err != nil {
				settings.Log.Fatal(err.Error())
			}

			stream, err := (*client).Upload(context.Background())

			if err != nil {
				helpers.ErrorHandler(settings.Log, err)
				return
			}

			err = stream.Send(&pb.UploadRequest{Data: &pb.UploadRequest_Metadata{
				Metadata: &pb.Metadata{
					FileName:        stat.Name(),
					FileContentType: mime.TypeByExtension(filepath.Ext(filePath)),
					FileSize:        stat.Size(),
					FileOffset:      remoteOffset,
				},
			}})

			if err != nil {
				settings.Log.Fatal(err.Error())
			}

			buf := make([]byte, 64*1024)

			for {
				n, err := file.Read(buf)

				if err == io.EOF {
					break
				}

				stream.Send(&pb.UploadRequest{Data: &pb.UploadRequest_Chunk{Chunk: buf[:n]}})
			}

			response, err := stream.CloseAndRecv()

			if err != nil {
				helpers.ErrorHandler(settings.Log, err)

				return
			}

			fmt.Println("✅ Файл загружен")
			fmt.Println("---")

			if response != nil && response.Secrets != nil && response.Secrets.Values != nil {
				helpers.PrintResult(response.Secrets.Values)
			}
		},
	}

	return cmd
}
