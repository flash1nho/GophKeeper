package actions

import (
	"context"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"

	"github.com/flash1nho/GophKeeper/app/client/helpers"
	"github.com/flash1nho/GophKeeper/config"
	pb "github.com/flash1nho/GophKeeper/internal/grpc"
	"github.com/spf13/cobra"
)

const (
	ChunkSize = 1024 * 1024 // 1MB
)

func SecretsUploadCommand(client *pb.GophKeeperPrivateServiceClient, settings config.SettingsObject) *cobra.Command {
	var path string

	cmd := &cobra.Command{
		Use:   "upload",
		Short: "Загрузка файла",
		Run: func(cmd *cobra.Command, args []string) {
			file, err := os.Open(path)

			if err != nil {
				settings.Log.Fatal(err.Error())
			}

			defer file.Close()

			stat, err := file.Stat()

			if err != nil {
				settings.Log.Fatal(err.Error())
			}

			statusResp, err := (*client).GetUploadStatus(context.Background(), &pb.UploadStatusRequest{
				FileName: stat.Name(),
			})

			if err != nil {
				helpers.ErrorHandler(settings.Log, err)

				return
			}

			var remoteOffset int64

			if statusResp != nil {
				remoteOffset = statusResp.FileOffset
			}

			totalSize := stat.Size()

			if remoteOffset >= totalSize {
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
					FileContentType: mime.TypeByExtension(filepath.Ext(path)),
					FileSize:        totalSize,
					FileOffset:      remoteOffset,
				},
			}})

			if err != nil {
				settings.Log.Fatal(err.Error())
			}

			buf := make([]byte, ChunkSize)
			currentSent := remoteOffset

			fmt.Printf("🚀 Начинаю загрузку: %s (всего %d байт)\n", stat.Name(), totalSize)

			for {
				n, err := file.Read(buf)

				if err == io.EOF {
					break
				}
				if err != nil {
					settings.Log.Fatal(err.Error())
				}

				err = stream.Send(&pb.UploadRequest{Data: &pb.UploadRequest_Chunk{Chunk: buf[:n]}})

				if err != nil {
					helpers.ErrorHandler(settings.Log, err)

					return
				}
				currentSent += int64(n)
				percentage := float64(currentSent) / float64(totalSize) * 100

				fmt.Printf("\r📤 Загрузка: %.2f%% [%d / %d байт]", percentage, currentSent, totalSize)
			}

			fmt.Println()

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

	cmd.Flags().StringVarP(&path, "path", "", path, "Путь для загрузки файла (обязательно)")
	_ = cmd.MarkFlagRequired("path")

	return cmd
}
