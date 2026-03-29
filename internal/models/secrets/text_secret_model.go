package secrets

import (
	"errors"

	"github.com/flash1nho/GophKeeper/config"
	"github.com/flash1nho/GophKeeper/internal/security"
)

var (
	ErrContentEmpty    = errors.New("'content' не может быть пустым")
	ErrUpdateDateEmpty = errors.New("укажите атрибут для обновления: 'content'")
)

type Text struct {
	BaseSecret

	Content string `json:"content"`
}

func NewText(userID int, settings config.SettingsObject) *Text {
	return &Text{
		BaseSecret: BaseSecret{
			UserID:        userID,
			CryptoManager: security.NewCryptoManager(settings.MasterKey),
		},
	}
}

func (t *Text) GetBaseSecret() *BaseSecret {
	return &t.BaseSecret
}

func (t *Text) GetType() string {
	return "Text"
}

func (t *Text) GetSecret() any {
	return t
}

func (t *Text) CreateValidate() error {
	if t.Content == "" {
		return ErrContentEmpty
	}

	return nil
}

func (t *Cred) UpdateValidate() error {
	if t.Content == "" {
		return ErrUpdateDateEmpty
	}

	return nil
}
