package secrets

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/flash1nho/GophKeeper/internal/security"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrInvalidCardData   = errors.New("недопустимые данные карты")
	ErrInvalidCardNumber = errors.New("недопустимый номер карты")
	ErrInvalidCardExpiry = errors.New("недопустимый срок действия карты")
	ErrInvalidCardHolder = errors.New("недопустимый владелец карты")
	ErrInvalidCardCVV    = errors.New("недопустимый CVV")
)

var numberRegex = regexp.MustCompile(`(\d{4})`)
var holderRegex = regexp.MustCompile(`^[a-zA-Z ]{2,30}$`)
var cvvRegex = regexp.MustCompile(`^[0-9]{3,4}$`)

type Card struct {
	BaseSecret

	CardType string `json:"card_type" ignore:"true"`
	Number   string `json:"number"`
	Expiry   string `json:"expiry"`
	Holder   string `json:"holder"`
	CVV      string `json:"cvv"`
}

func NewCard(userID int, masterKey []byte, pool *pgxpool.Pool) *Card {
	return &Card{
		BaseSecret: BaseSecret{
			UserID:        userID,
			CryptoManager: security.NewCryptoManager(masterKey),
			pool:          pool,
		},
	}
}

func (s *Card) GetType() string {
	return "Card"
}

func (s *Card) GetSecret() any {
	return s
}

func (s *Card) CreateValidate(ctx context.Context) error {
	fields := s.validationFields()

	for _, f := range fields {
		if strings.TrimSpace(f.value) == "" {
			return fmt.Errorf("%w: поле '%s' не может быть пустым", ErrInvalidCardData, f.name)
		} else if f.name == "number" && !validateNumber(f.value) {
			return ErrInvalidCardNumber
		} else if f.name == "expiry" && !validateExpiry(f.value) {
			return ErrInvalidCardExpiry
		} else if f.name == "holder" && !validateHolder(f.value) {
			return ErrInvalidCardHolder
		} else if f.name == "cvv" && !validateCVV(f.value) {
			return ErrInvalidCardCVV
		}
	}

	return nil
}

func (s *Card) UpdateValidate() error {
	if s.Number == "" && s.Expiry == "" && s.Holder == "" && s.CVV == "" {
		return errors.New("нужно указать хотя бы один атрибут для обновления: 'number', 'expiry', 'holder', 'cvv'")
	}

	fields := s.validationFields()

	for _, f := range fields {
		if strings.TrimSpace(f.value) != "" {
			if f.name == "number" && !validateNumber(f.value) {
				return ErrInvalidCardNumber
			} else if f.name == "expiry" && !validateExpiry(f.value) {
				return ErrInvalidCardExpiry
			} else if f.name == "holder" && !validateHolder(f.value) {
				return ErrInvalidCardHolder
			} else if f.name == "cvv" && !validateCVV(f.value) {
				return ErrInvalidCardCVV
			}
		}
	}

	return nil
}

func (s *Card) FileExists(ctx context.Context) (bool, error) {
	return false, ErrNotImplemented
}

func (s *Card) validationFields() []struct{ name, value string } {
	return []struct {
		name  string
		value string
	}{
		{"number", s.Number},
		{"expiry", s.Expiry},
		{"holder", s.Holder},
		{"cvv", s.CVV},
	}
}

func validateNumber(number string) bool {
	number = strings.ReplaceAll(number, " ", "")

	if len(number) < 13 || len(number) > 19 {
		return false
	}

	var sum int
	alternate := false

	for i := len(number) - 1; i >= 0; i-- {
		n, _ := strconv.Atoi(string(number[i]))

		if alternate {
			n *= 2

			if n > 9 {
				n -= 9
			}
		}

		sum += n
		alternate = !alternate
	}

	return sum%10 == 0
}

func validateExpiry(exp string) bool {
	t, err := time.Parse("01/06", exp)

	if err != nil {
		return false
	}

	now := time.Now()
	lastDay := t.AddDate(0, 1, 0).Add(-time.Second)

	return lastDay.After(now)
}

func validateHolder(holder string) bool {
	return holderRegex.MatchString(holder)
}

func validateCVV(cvv string) bool {
	return cvvRegex.MatchString(cvv)
}

func GetCardType(number string) string {
	patterns := map[string]string{
		"Visa":       "^4",
		"MasterCard": "^5[1-5]",
		"Amex":       "^3[47]",
		"Mir":        "^22",
	}

	for name, pattern := range patterns {
		if match, _ := regexp.MatchString(pattern, number); match {
			return name
		}
	}

	return "Unknown"
}

func FormatCardNumber(number string) string {
	number = strings.ReplaceAll(number, " ", "")
	formatted := numberRegex.ReplaceAllString(number, "$1 ")

	return strings.TrimSpace(formatted)
}
