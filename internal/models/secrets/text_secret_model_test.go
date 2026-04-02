package secrets

import (
	"context"
	"testing"

	"github.com/flash1nho/GophKeeper/internal/security"
	"github.com/stretchr/testify/assert"
)

func TestTextCreateValidate(t *testing.T) {
	masterKey := []byte("test_master_key_for_testing_purposes")
	cryptoManager := security.NewCryptoManager(masterKey)

	tests := []struct {
		name      string
		text      *Text
		shouldErr bool
	}{
		{
			name: "Valid text content",
			text: &Text{
				BaseSecret: BaseSecret{CryptoManager: cryptoManager},
				Content:    "This is important text data",
			},
			shouldErr: false,
		},
		{
			name: "Empty content",
			text: &Text{
				BaseSecret: BaseSecret{CryptoManager: cryptoManager},
				Content:    "",
			},
			shouldErr: true,
		},
		{
			name: "Whitespace only content",
			text: &Text{
				BaseSecret: BaseSecret{CryptoManager: cryptoManager},
				Content:    "   \n\t  ",
			},
			shouldErr: true,
		},
		{
			name: "Long text content",
			text: &Text{
				BaseSecret: BaseSecret{CryptoManager: cryptoManager},
				Content:    string(make([]byte, 10000)),
			},
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.text.CreateValidate(context.Background())
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTextUpdateValidate(t *testing.T) {
	masterKey := []byte("test_master_key_for_testing_purposes")
	cryptoManager := security.NewCryptoManager(masterKey)

	tests := []struct {
		name      string
		text      *Text
		shouldErr bool
	}{
		{
			name: "Valid content update",
			text: &Text{
				BaseSecret: BaseSecret{CryptoManager: cryptoManager},
				Content:    "updated content",
			},
			shouldErr: false,
		},
		{
			name: "Empty content update",
			text: &Text{
				BaseSecret: BaseSecret{CryptoManager: cryptoManager},
				Content:    "",
			},
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.text.UpdateValidate()
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTextGetType(t *testing.T) {
	text := &Text{}
	assert.Equal(t, "Text", text.GetType())
}

func TestTextGetSecret(t *testing.T) {
	masterKey := []byte("test_master_key")
	cryptoManager := security.NewCryptoManager(masterKey)
	text := &Text{
		BaseSecret: BaseSecret{CryptoManager: cryptoManager},
		Content:    "secret content",
	}

	secret := text.GetSecret()
	assert.Equal(t, text, secret)
}

func TestTextFileExists(t *testing.T) {
	masterKey := []byte("test_master_key")
	cryptoManager := security.NewCryptoManager(masterKey)
	text := &Text{
		BaseSecret: BaseSecret{CryptoManager: cryptoManager},
	}

	exists, err := text.FileExists(context.Background())
	assert.Error(t, err)
	assert.False(t, exists)
	assert.Equal(t, ErrNotImplemented, err)
}
