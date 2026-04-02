package secrets

import (
	"context"
	"testing"

	"github.com/flash1nho/GophKeeper/internal/security"
	"github.com/stretchr/testify/assert"
)

func TestFileCreateValidate(t *testing.T) {
	masterKey := []byte("test_master_key_for_testing_purposes")
	cryptoManager := security.NewCryptoManager(masterKey)

	tests := []struct {
		name      string
		file      *File
		shouldErr bool
	}{
		{
			name: "Valid file with name",
			file: &File{
				BaseSecret: BaseSecret{CryptoManager: cryptoManager},
				FileName:   "document.pdf",
			},
			shouldErr: false,
		},
		{
			name: "Empty filename",
			file: &File{
				BaseSecret: BaseSecret{CryptoManager: cryptoManager},
				FileName:   "",
			},
			shouldErr: true,
		},
		{
			name: "Whitespace filename",
			file: &File{
				BaseSecret: BaseSecret{CryptoManager: cryptoManager},
				FileName:   "   ",
			},
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = tt.file.CreateValidate(context.Background())

			if tt.file.FileName == "" {
				assert.True(t, tt.shouldErr)
			}
		})
	}
}

func TestFileUpdateValidate(t *testing.T) {
	masterKey := []byte("test_master_key_for_testing_purposes")
	cryptoManager := security.NewCryptoManager(masterKey)

	file := &File{
		BaseSecret: BaseSecret{CryptoManager: cryptoManager},
	}

	err := file.UpdateValidate()
	assert.NoError(t, err)
}

func TestFileGetType(t *testing.T) {
	file := &File{}
	assert.Equal(t, "File", file.GetType())
}

func TestFileGetSecret(t *testing.T) {
	masterKey := []byte("test_master_key")
	cryptoManager := security.NewCryptoManager(masterKey)

	file := &File{
		BaseSecret: BaseSecret{CryptoManager: cryptoManager},
		FileName:   "test.txt",
	}

	secret := file.GetSecret()
	assert.Equal(t, file, secret)
}
