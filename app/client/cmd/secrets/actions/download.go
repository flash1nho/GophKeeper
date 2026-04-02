package actions

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/flash1nho/GophKeeper/app/client/helpers"
	"github.com/flash1nho/GophKeeper/config"
	pb "github.com/flash1nho/GophKeeper/internal/grpc"
	"github.com/spf13/cobra"
)

func SecretsDownloadCommand(client *pb.GophKeeperPrivateServiceClient, settings config.SettingsObject) *cobra.Command {
	var id int32
	var outputPath string

	cmd := &cobra.Command{
		Use:   "download",
		Short: "Скачивание файла с докачкой",
		Run: func(cmd *cobra.Command, args []string) {
			var localOffset int64 = 0

			stat, err := os.Stat(outputPath)

			if err == nil {
				localOffset = stat.Size()
				fmt.Printf("🔄 Найден локальный фрагмент (%d байт), запрашиваю докачку...\n", localOffset)
			}

			stream, err := (*client).Download(context.Background(), &pb.DownloadRequest{
				ID:         id,
				FileOffset: localOffset,
			})

			if err != nil {
				helpers.ErrorHandler(settings.Log, err)
				return
			}

			file, err := os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

			if err != nil {
				settings.Log.Fatal(fmt.Sprintf("не удалось открыть файл для записи: %v", err))
			}

			defer file.Close()

			fmt.Println("🚀 Начинаю загрузку...")

			var totalDownloaded int64 = localOffset
			var totalSize int64

			for {
				resp, err := stream.Recv()

				if err == io.EOF {
					break
				}
				if err != nil {
					helpers.ErrorHandler(settings.Log, err)

					return
				}

				n, err := file.Write(resp.Chunk)

				if err != nil {
					settings.Log.Fatal(fmt.Sprintf("ошибка при записи в файл: %v", err))
				}

				totalDownloaded += int64(n)

				if totalSize == 0 && resp.FileOffset > 0 {
					totalSize = resp.FileOffset
				}

				if totalSize > 0 {
					percentage := float64(totalDownloaded) / float64(totalSize) * 100
					fmt.Printf("\r📥 Прогресс: %.2f%% (%d / %d байт)", percentage, totalDownloaded, totalSize)
				} else {
					fmt.Printf("\r📥 Скачано: %d байт", totalDownloaded)
				}
			}

			fmt.Println("\n✅ Файл успешно скачан")
		},
	}

	cmd.Flags().Int32VarP(&id, "id", "", id, "ID секрета (обязательно)")
	cmd.Flags().StringVarP(&outputPath, "out", "", outputPath, "Путь для сохранения файла (обязательно)")
	cmd.MarkFlagRequired("id")
	cmd.MarkFlagRequired("out")

	return cmd
}
