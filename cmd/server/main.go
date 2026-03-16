package main

import (
	"fmt"

	"github.com/flash1nho/GophKeeper/internal/config"
	"github.com/flash1nho/GophKeeper/internal/facade"
	"github.com/flash1nho/GophKeeper/internal/grpc"
	"github.com/flash1nho/GophKeeper/internal/handler"
	"github.com/flash1nho/GophKeeper/internal/service"
)

func main() {
	settings := config.Settings()
	// store, err := storage.NewStorage(settings.FilePath, settings.DatabaseDSN)

	if err != nil {
		settings.Log.Error(fmt.Sprint(err))
	}

	f := facade.NewFacade(store, settings.Server2.BaseURL)
	h := handler.NewHandler(f, settings)
	gh := grpc.NewHandler(f)
	service.NewService(h, gh, settings).Run()
}
