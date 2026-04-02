package secrets

import (
	"context"
	"testing"

	"github.com/flash1nho/GophKeeper/internal/security"
	"github.com/stretchr/testify/assert"
)

func TestCredCreateValidate(t *testing.T) {
	masterKey := []byte("test_master_key_for_testing_purposes")
	cryptoManager := security.NewCryptoManager(masterKey)

	tests := []struct {
		name      string
		cred      *Cred
		shouldErr bool
	}{
		{
			name: "Valid credentials",
			cred: &Cred{
				BaseSecret: BaseSecret{CryptoManager: cryptoManager},
				Login:      "testuser",
				Password:   "password123",
			},
			shouldErr: false,
		},
		{
			name: "Empty login",
			cred: &Cred{
				BaseSecret: BaseSecret{CryptoManager: cryptoManager},
				Password:   "password123",
			},
			shouldErr: true,
		},
		{
			name: "Empty password",
			cred: &Cred{
				BaseSecret: BaseSecret{CryptoManager: cryptoManager},
				Login:      "testuser",
			},
			shouldErr: true,
		},
		{
			name: "Both empty",
			cred: &Cred{
				BaseSecret: BaseSecret{CryptoManager: cryptoManager},
			},
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cred.CreateValidate(context.Background())
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCredUpdateValidate(t *testing.T) {
	masterKey := []byte("test_master_key_for_testing_purposes")
	cryptoManager := security.NewCryptoManager(masterKey)

	tests := []struct {
		name      string
		cred      *Cred
		shouldErr bool
	}{
		{
			name: "Update login only",
			cred: &Cred{
				BaseSecret: BaseSecret{CryptoManager: cryptoManager},
				Login:      "newlogin",
			},
			shouldErr: false,
		},
		{
			name: "Update password only",
			cred: &Cred{
				BaseSecret: BaseSecret{CryptoManager: cryptoManager},
				Password:   "newpassword",
			},
			shouldErr: false,
		},
		{
			name: "Empty update",
			cred: &Cred{
				BaseSecret: BaseSecret{CryptoManager: cryptoManager},
			},
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cred.UpdateValidate()
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
