package secrets

import (
	"errors"
	"fmt"
	"strings"

	"github.com/flash1nho/GophKeeper/internal/security"
)

var ErrInvalidTextData = errors.New("недопустимые текстовые данные")

type Text struct {
	BaseSecret

	Content string `json:"content"`
}

func NewText(userID int, masterKey []byte) *Text {
	return &Text{
		BaseSecret: BaseSecret{
			UserID:        userID,
			CryptoManager: security.NewCryptoManager(masterKey),
		},
	}
}

func (s *Text) GetType() string {
	return "Text"
}

func (s *Text) GetSecret() any {
	return s
}

func (s *Text) CreateValidate() error {
	if strings.TrimSpace(s.Content) == "" {
		return fmt.Errorf("%w: поле 'content' не может быть пустым", ErrInvalidTextData)
	}

	return nil
}

func (s *Text) UpdateValidate() error {
	if strings.TrimSpace(s.Content) == "" {
		return errors.New("нужно указать атрибут для обновления: 'content'")
	}

	return nil
}
