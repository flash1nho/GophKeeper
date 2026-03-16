package facade

import (
	"context"

	"github.com/flash1nho/GophKeeper/internal/authenticator"
	"golang.org/x/mod/sumdb/storage"
)

type Facade struct {
	Store   *storage.Storage
	BaseURL string
}

func NewFacade(store *storage.Storage, BaseURL string) *Facade {
	return &Facade{
		Store:   store,
		BaseURL: BaseURL,
	}
}

func (f *Facade) GetUserFromContext(ctx context.Context) (string, error) {
	userID, err := authenticator.FromContext(ctx)

	if err != nil {
		return "", nil
	}

	return userID, nil
}
