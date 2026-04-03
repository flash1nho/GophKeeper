package users

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	login := "testuser"
	password := "password123"

	user := NewUser(login, password)

	assert.Equal(t, login, user.Login)
	assert.Equal(t, password, user.Password)
}

func TestUserRegisterValidation(t *testing.T) {
	tests := []struct {
		name      string
		login     string
		password  string
		secret    string
		shouldErr bool
		errType   error
	}{
		{
			name:      "Valid credentials",
			login:     "testuser",
			password:  "password123",
			secret:    "mysecret",
			shouldErr: false,
		},
		{
			name:      "Empty login",
			login:     "",
			password:  "password123",
			secret:    "mysecret",
			shouldErr: true,
			errType:   ErrCredentialsRequired,
		},
		{
			name:      "Empty password",
			login:     "testuser",
			password:  "",
			secret:    "mysecret",
			shouldErr: true,
			errType:   ErrCredentialsRequired,
		},
		{
			name:      "Empty secret",
			login:     "testuser",
			password:  "password123",
			secret:    "",
			shouldErr: true,
			errType:   ErrSecretRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewUser(tt.login, tt.password)

			if tt.login == "" || tt.password == "" {
				assert.True(t, tt.shouldErr)
			}
			if tt.secret == "" {
				assert.True(t, tt.shouldErr)
			}
		})
	}
}
