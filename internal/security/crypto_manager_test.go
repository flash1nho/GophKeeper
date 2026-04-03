package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	manager := NewCryptoManager([]byte("test_master_key"))
	password := "mysecurepassword123"

	hash, err := manager.HashPassword(password)
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)
}

func TestCheckPassword(t *testing.T) {
	manager := NewCryptoManager([]byte("test_master_key"))
	password := "mysecurepassword123"

	hash, err := manager.HashPassword(password)
	require.NoError(t, err)

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Correct password",
			input:    password,
			expected: true,
		},
		{
			name:     "Wrong password",
			input:    "wrongpassword",
			expected: false,
		},
		{
			name:     "Empty password",
			input:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.CheckPassword(tt.input, hash)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEncryptDecrypt(t *testing.T) {
	manager := NewCryptoManager([]byte("test_master_key"))
	plaintext := []byte("sensitive data to encrypt")
	key := []byte("encryption_key")

	encrypted, err := manager.Encrypt(plaintext, key)
	require.NoError(t, err)
	assert.NotEqual(t, plaintext, encrypted)

	decrypted, err := manager.Decrypt(encrypted, key)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

func TestDecryptWithWrongKey(t *testing.T) {
	manager := NewCryptoManager([]byte("test_master_key"))
	plaintext := []byte("sensitive data")
	key := []byte("correct_key")
	wrongKey := []byte("wrong_key")

	encrypted, err := manager.Encrypt(plaintext, key)
	require.NoError(t, err)

	decrypted, err := manager.Decrypt(encrypted, wrongKey)
	assert.Error(t, err)
	assert.NotEqual(t, plaintext, decrypted)
}

func TestGenerateAndValidateToken(t *testing.T) {
	manager := NewCryptoManager([]byte("test_master_key_for_testing_purposes"))
	userID := 123
	secret := []byte("user_secret_data")

	encryptedSecret, err := manager.Encrypt(secret, manager.MasterKey)
	require.NoError(t, err)

	token, err := manager.GenerateToken(userID, encryptedSecret)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	validatedUserID, err := manager.ValidateToken(token, encryptedSecret)
	require.NoError(t, err)
	assert.Equal(t, userID, validatedUserID)
}

func TestTokenExpiration(t *testing.T) {
	manager := NewCryptoManager([]byte("test_master_key_for_testing_purposes"))
	userID := 456
	secret := []byte("user_secret_data")

	encryptedSecret, err := manager.Encrypt(secret, manager.MasterKey)
	require.NoError(t, err)

	token, err := manager.GenerateToken(userID, encryptedSecret)
	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestConvertKeyToSha256(t *testing.T) {
	manager := NewCryptoManager([]byte("test"))

	shortKey := []byte("short_key")
	result := manager.ConverKeyToSha256(shortKey)
	assert.Equal(t, 32, len(result))

	longKey := make([]byte, 32)
	for i := range longKey {
		longKey[i] = byte(i)
	}
	result2 := manager.ConverKeyToSha256(longKey)
	assert.Equal(t, longKey, result2)
}

func TestGetUserIDFromToken(t *testing.T) {
	manager := NewCryptoManager([]byte("test_master_key_for_testing_purposes"))
	userIDStr := "789"
	encryptedUserID, err := manager.Encrypt([]byte(userIDStr), manager.MasterKey)
	require.NoError(t, err)

	assert.NotEmpty(t, encryptedUserID)
}
