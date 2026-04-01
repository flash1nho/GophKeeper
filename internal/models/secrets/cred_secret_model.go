package secrets

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/flash1nho/GophKeeper/internal/security"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrInvalidCredData = errors.New("недопустимые учетные данные")

type Cred struct {
	BaseSecret

	Login    string `json:"login"`
	Password string `json:"password"`
}

func NewCred(userID int, masterKey []byte, pool *pgxpool.Pool) *Cred {
	return &Cred{
		BaseSecret: BaseSecret{
			UserID:        userID,
			CryptoManager: security.NewCryptoManager(masterKey),
			pool:          pool,
		},
	}
}

func (s *Cred) GetType() string {
	return "Cred"
}

func (s *Cred) GetSecret() any {
	return s
}

func (s *Cred) CreateValidate(ctx context.Context) error {
	fields := []struct {
		name  string
		value string
	}{
		{"login", s.Login},
		{"password", s.Password},
	}

	for _, f := range fields {
		if strings.TrimSpace(f.value) == "" {
			return fmt.Errorf("%w: поле '%s' не может быть пустым", ErrInvalidCredData, f.name)
		}
	}

	return nil
}

func (s *Cred) UpdateValidate() error {
	if strings.TrimSpace(s.Login) == "" && strings.TrimSpace(s.Password) == "" {
		return errors.New("нужно указать хотя бы один атрибут для обновления: 'login', 'password'")
	}

	return nil
}

func (s *Cred) FileExists(ctx context.Context) (bool, error) {
	return false, ErrNotImplemented
}

func (s *Cred) GetFileOffset() (int64, error) {
	return 0, ErrNotImplemented
}

func (s *Cred) SetFileOffset(fileOffset int64) {
}
