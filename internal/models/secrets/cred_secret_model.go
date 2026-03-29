package secrets

import (
	"errors"

	"github.com/flash1nho/GophKeeper/config"
	"github.com/flash1nho/GophKeeper/internal/security"
)

var (
	ErrLoginEmpty      = errors.New("'login' не может быть пустым")
	ErrPasswordEmpty   = errors.New("'login' не может быть пустым")
	ErrUpdateCredEmpty = errors.New("укажите хотя бы одно из доступных атрибутов для обновления: 'login' или 'password'")
)

type Cred struct {
	BaseSecret

	Login    string `json:"login"`
	Password string `json:"password"`
}

func NewCred(userID int, settings config.SettingsObject) *Cred {
	return &Cred{
		BaseSecret: BaseSecret{
			UserID:        userID,
			CryptoManager: security.NewCryptoManager(settings.MasterKey),
		},
	}
}

func (t *Cred) GetBaseSecret() *BaseSecret {
	return &t.BaseSecret
}

func (t *Cred) GetType() string {
	return "Cred"
}

func (t *Cred) GetSecret() any {
	return t
}

func (t *Cred) CreateValidate() error {
	if t.Login == "" {
		return ErrLoginEmpty
	}

	if t.Password == "" {
		return ErrPasswordEmpty
	}

	return nil
}

func (t *Cred) UpdateValidate() error {
	if t.Login == "" && t.Password == "" {
		return ErrUpdateCredEmpty
	}

	return nil
}
