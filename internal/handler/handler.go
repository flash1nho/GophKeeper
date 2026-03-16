package handler

import (
	"github.com/flash1nho/go-musthave-shortener-tpl/internal/config"
	"github.com/flash1nho/go-musthave-shortener-tpl/internal/facade"
)

func NewHandler(facade *facade.Facade, settings config.SettingsObject) *Handler {
	return &Handler{
		Facade: facade,
		log:    settings.Log,
	}
}
