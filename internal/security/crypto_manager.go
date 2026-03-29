package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidData       = errors.New("недопустимые зашифрованные данные")
	ErrPadding           = errors.New("недопустимое дополнение")
	ErrTokenExpired      = errors.New("срок действия токена истек")
	ErrTokenInvalid      = errors.New("невалидный токен")
	ErrDecryptUserKey    = errors.New("ошибка расшифровки ключа пользователя")
	ErrUnexpectedSigning = errors.New("неожиданный метод подписи")
)

type CryptoManager struct {
	masterKey []byte
}

func NewCryptoManager(masterKey []byte) *CryptoManager {
	key := sha256.Sum256(masterKey)

	return &CryptoManager{masterKey: key[:]}
}

func (m *CryptoManager) Encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(m.masterKey)

	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)

	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func (m *CryptoManager) Decrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(m.masterKey)

	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)

	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()

	if len(data) < nonceSize {
		return nil, ErrInvalidData
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	return gcm.Open(nil, nonce, ciphertext, nil)
}

func (m *CryptoManager) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(hash), err
}

func (m *CryptoManager) CheckPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func (m *CryptoManager) GenerateToken(userID int, encryptedSecret []byte) (string, error) {
	userKey, err := m.Decrypt(encryptedSecret)

	if err != nil {
		return "", ErrDecryptUserKey
	}

	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(userKey)
}

func (m *CryptoManager) ValidateToken(tokenStr string, encryptedSecret []byte) (int, error) {
	userKey, err := m.Decrypt(encryptedSecret)

	if err != nil {
		return 0, err
	}

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return 0, ErrUnexpectedSigning
		}

		return userKey, nil
	})

	if err != nil {
		var vErr *jwt.ValidationError

		if errors.As(err, &vErr) {
			if vErr.Errors&jwt.ValidationErrorExpired != 0 {
				return 0, ErrTokenExpired
			}
		}

		return 0, ErrTokenInvalid
	}

	if !token.Valid {
		return 0, ErrTokenInvalid
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if sub, ok := claims["sub"].(float64); ok {
			return int(sub), nil
		}
	}

	return 0, ErrTokenInvalid
}
