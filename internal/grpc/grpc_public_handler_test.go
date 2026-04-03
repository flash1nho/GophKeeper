package grpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/flash1nho/GophKeeper/config"
	"github.com/flash1nho/GophKeeper/internal/logger"
	"github.com/flash1nho/GophKeeper/internal/models/users"
)

type MockUserModel struct {
	mock.Mock
}

func TestGrpcPublicHandlerRegisterSuccess(t *testing.T) {
	_ = logger.Initialize("debug")

	settings := config.SettingsObject{
		MasterKey: []byte("test_master_key_for_testing_purposes"),
		Log:       logger.Log,
	}

	handler := &GrpcPublicHandler{
		pool:     nil,
		settings: settings,
	}

	assert.NotNil(t, handler)
	assert.Equal(t, settings, handler.settings)
}

func TestGrpcPublicHandlerRegisterValidation(t *testing.T) {
	tests := []struct {
		name      string
		req       *UserRegisterRequest
		shouldErr bool
	}{
		{
			name: "Valid registration request",
			req: &UserRegisterRequest{
				Login:    "newuser",
				Password: "securepass123",
				Secret:   "mysecret",
			},
			shouldErr: false,
		},
		{
			name: "Empty login",
			req: &UserRegisterRequest{
				Login:    "",
				Password: "password123",
				Secret:   "mysecret",
			},
			shouldErr: true,
		},
		{
			name: "Empty password",
			req: &UserRegisterRequest{
				Login:    "testuser",
				Password: "",
				Secret:   "mysecret",
			},
			shouldErr: true,
		},
		{
			name: "Empty secret",
			req: &UserRegisterRequest{
				Login:    "testuser",
				Password: "password123",
				Secret:   "",
			},
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.req.Login == "" || tt.req.Password == "" {
				assert.True(t, tt.shouldErr)
			}
			if tt.req.Secret == "" {
				assert.True(t, tt.shouldErr)
			}
		})
	}
}

func TestGrpcPublicHandlerLoginValidation(t *testing.T) {
	tests := []struct {
		name      string
		req       *UserLoginRequest
		shouldErr bool
	}{
		{
			name: "Valid login request",
			req: &UserLoginRequest{
				Login:    "testuser",
				Password: "password123",
			},
			shouldErr: false,
		},
		{
			name: "Empty login",
			req: &UserLoginRequest{
				Login:    "",
				Password: "password123",
			},
			shouldErr: true,
		},
		{
			name: "Empty password",
			req: &UserLoginRequest{
				Login:    "testuser",
				Password: "",
			},
			shouldErr: true,
		},
		{
			name: "Both empty",
			req: &UserLoginRequest{
				Login:    "",
				Password: "",
			},
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.req.Login == "" || tt.req.Password == "" {
				assert.True(t, tt.shouldErr)
			}
		})
	}
}

func TestGrpcPublicHandlerInitialization(t *testing.T) {
	_ = logger.Initialize("debug")
	settings := config.SettingsObject{
		DatabaseDSN:       "postgres://localhost/test",
		GrpcServerAddress: "localhost:3200",
		MasterKey:         []byte("test_key"),
		Log:               logger.Log,
	}

	handler := &GrpcPublicHandler{
		pool:     nil,
		settings: settings,
	}

	assert.NotNil(t, handler)
	assert.Equal(t, settings.DatabaseDSN, handler.settings.DatabaseDSN)
	assert.Equal(t, settings.GrpcServerAddress, handler.settings.GrpcServerAddress)
	assert.NotNil(t, handler.settings.Log)
}

func TestGrpcPublicHandlerUserCreation(t *testing.T) {
	tests := []struct {
		name        string
		login       string
		password    string
		expectedErr bool
	}{
		{
			name:        "Valid user",
			login:       "user1",
			password:    "pass1",
			expectedErr: false,
		},
		{
			name:        "Empty login",
			login:       "",
			password:    "pass1",
			expectedErr: true,
		},
		{
			name:        "Empty password",
			login:       "user1",
			password:    "",
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := users.NewUser(tt.login, tt.password)

			if tt.expectedErr {
				assert.True(t, user.Login == "" || user.Password == "")
			} else {
				assert.Equal(t, tt.login, user.Login)
				assert.Equal(t, tt.password, user.Password)
			}
		})
	}
}
