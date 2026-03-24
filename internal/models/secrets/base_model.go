package secrets

import (
	"encoding/json"
	"time"

	"github.com/Masterminds/squirrel"
)

type BaseRecord struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"default:current_timestamp"`
}

type BankCard struct {
	BaseRecord

	Type           string `gorm:"default:BankCard"`
	CardMask       string `gorm:"column:card_mask"`
	CardHolder     string `gorm:"column:card_holder"`
	CardCVV        string `gorm:"column:card_cvv"`
	CardExpiryDate string `gorm:"column:card_expiry_date"`
}

type FileUpload struct {
	BaseRecord

	Type        string `gorm:"default:FileUpload"`
	FileName    string `gorm:"column:file_name"`
	ContentType string `gorm:"column:content_type"`
	FileSize    int    `gorm:"column:file_size"`
}

type TextNote struct {
	BaseRecord

	Type    string `gorm:"default:TextNote"`
	Content string `gorm:"default:''"`
}

type Processor[T BankCard | FileUpload | TextNote] struct{}

func (p Processor[T]) Create(obj T, userID int) (*int, error) {
	properties, err := json.Marshal(obj)

	if err != nil {
		return nil, err
	}

	sql, args, err := squirrel.Insert("secrets").
		Columns("user_id", "properties", "created_at").
		Values(userID, properties, time.Now().UTC()).
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	err = pool.QueryRow(ctx, sql, args...).Scan(&obj.ID)

	if err != nil {
		return nil, err
	}

	return obj.ID, nil
}
