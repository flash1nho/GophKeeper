package main

import (
	"fmt"

	"github.com/flash1nho/GophKeeper/internal/config"
	"github.com/flash1nho/GophKeeper/internal/db"
	"github.com/flash1nho/GophKeeper/internal/grpc"
	"github.com/flash1nho/GophKeeper/internal/service"
)

func main() {
	settings := config.Settings()
	pool, err := db.NewDB(settings.DatabaseURI)

	if err != nil {
		settings.Log.Error(fmt.Sprint(err))
	}

	gh := grpc.NewHandler(pool)
	service.NewService(gh, settings).Run()
}
