package facade

import (
	"context"

	"github.com/flash1nho/GophKeeper/internal/interceptors"
)

type Facade struct{}

func NewFacade() *Facade {
	return &Facade{}
}

func (f *Facade) GetUserFromContext(ctx context.Context) (int, error) {
	userID, err := interceptors.GetUserFromContext(ctx)

	if err != nil {
		return 0, err
	}

	return userID, nil
}
