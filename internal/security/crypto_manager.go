package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"strconv"
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
	MasterKey []byte
}

func NewCryptoManager(masterKey []byte) *CryptoManager {
	return &CryptoManager{MasterKey: masterKey}
}

func (m *CryptoManager) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(hash), err
}

func (m *CryptoManager) CheckPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func (m *CryptoManager) Encrypt(data []byte, key []byte) ([]byte, error) {
	key = m.converKeyToSha256(key)

	block, err := aes.NewCipher(key)

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

	return gcm.Seal(nonce, nonce, data, nil), nil
}

func (m *CryptoManager) Decrypt(data []byte, key []byte) ([]byte, error) {
	key = m.converKeyToSha256(key)

	block, err := aes.NewCipher(key)

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

func (m *CryptoManager) GenerateToken(userID int, enctryptedSecret []byte) (string, error) {
	secret, err := m.Decrypt(enctryptedSecret, m.MasterKey)

	if err != nil {
		return "", err
	}

	strUserID := strconv.Itoa(userID)
	encryptedSub, err := m.Encrypt([]byte(strUserID), m.MasterKey)

	if err != nil {
		return "", err
	}

	encodedSub := base64.StdEncoding.EncodeToString(encryptedSub)

	claims := jwt.MapClaims{
		"sub": encodedSub,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secret)
}

func (m *CryptoManager) ValidateToken(tokenStr string, enctryptedSecret []byte) (int, error) {
	secret, err := m.Decrypt(enctryptedSecret, m.MasterKey)

	if err != nil {
		return 0, err
	}

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigning
		}

		return secret, nil
	})

	if err != nil || !token.Valid {
		return 0, ErrTokenInvalid
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return 0, ErrTokenInvalid
	}

	return m.GetUserIDFromToken(claims)
}

func (m *CryptoManager) GetUserIDFromToken(claims jwt.MapClaims) (int, error) {
	encodedSub, ok := claims["sub"].(string)

	if !ok {
		return 0, ErrTokenInvalid
	}

	encryptedSub, err := base64.StdEncoding.DecodeString(encodedSub)

	if err != nil {
		return 0, ErrTokenInvalid
	}

	decryptedSub, err := m.Decrypt(encryptedSub, m.MasterKey)

	if err != nil {
		return 0, ErrDecryptUserKey
	}

	subStr := string(decryptedSub)
	userID, err := strconv.Atoi(subStr)

	if err != nil {
		return 0, ErrTokenInvalid
	}

	return userID, nil
}

func (m *CryptoManager) converKeyToSha256(key []byte) []byte {
	if len(key) != 32 {
		h := sha256.Sum256(key)
		key = h[:]
	}

	return key
}
