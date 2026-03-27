package secrets

import (
	"fmt"
)

func NewTextNote(userID int, content string) *TextNote {
	return &TextNote{
		BaseSecret: BaseSecret{UserID: userID},
		Content:    content,
	}
}

func (tn *TextNote) GetBaseSecret() *BaseSecret {
	return &tn.BaseSecret
}

func (tn *TextNote) GetType() string {
	return "TextNote"
}

func (tn *TextNote) GetPayload() any {
	return map[string]string{"content": tn.Content}
}

func (tn *TextNote) Validate() error {
	if tn.Content == "" {
		return fmt.Errorf("Content не может быть пустынм")
	}

	return nil
}
