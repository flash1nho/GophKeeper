package secrets

import "fmt"

func NewTextNote(content *string) {
	return &TextNote{
		Content: content,
	}
}

func (obj *TextNote) create(userID int) (*int, error) {
	if obj.Content == "" {
		return nil, fmt.Errorf("Content не может быть пустым")
	}

	objID, err := Create(obj, userID)

	if err != nil {
		return nil, err
	}

	return obj.ID, nil
}
