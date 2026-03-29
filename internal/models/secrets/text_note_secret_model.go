package secrets

import (
	"errors"
)

var (
	ErrContentEmpty = errors.New("текст не может быть пустым")
)

type Text struct {
	BaseSecret

	Content string
}

func NewText(userID int, content string) *Text {
	return &Text{
		BaseSecret: BaseSecret{UserID: userID},
		Content:    content,
	}
}

func (tn *Text) GetBaseSecret() *BaseSecret {
	return &tn.BaseSecret
}

func (tn *Text) GetType() string {
	return "Text"
}

func (tn *Text) GetPayload() any {
	return map[string]string{"content": tn.Content}
}

func (tn *Text) Validate() error {
	if tn.Content == "" {
		return ErrContentEmpty
	}

	return nil
}
