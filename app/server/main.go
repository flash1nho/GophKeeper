package main

import (
	"github.com/flash1nho/GophKeeper/config"
	"github.com/flash1nho/GophKeeper/internal/db"
	"github.com/flash1nho/GophKeeper/internal/facade"
	"github.com/flash1nho/GophKeeper/internal/grpc"
	"github.com/flash1nho/GophKeeper/internal/service"
)

func main() {
	settings := config.Settings()
	pool, err := db.NewDB(settings.DatabaseDSN)

	if err != nil {
		settings.Log.Fatal(err.Error())
	}

	f := facade.NewFacade()
	handler := grpc.NewGrpcHandler(pool, settings, f)
	service.NewService(handler).Run()
}
