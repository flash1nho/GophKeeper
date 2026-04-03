package interceptors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func TestGetUserIDFromContext(t *testing.T) {
	tests := []struct {
		name      string
		ctx       context.Context
		expected  int
		shouldErr bool
	}{
		{
			name:      "Valid user ID in context",
			ctx:       context.WithValue(context.Background(), userKey, 123),
			expected:  123,
			shouldErr: false,
		},
		{
			name:      "Missing user ID in context",
			ctx:       context.Background(),
			expected:  0,
			shouldErr: true,
		},
		{
			name:      "Wrong type in context",
			ctx:       context.WithValue(context.Background(), userKey, "not_an_int"),
			expected:  0,
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := GetUserIDFromContext(tt.ctx)
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, userID)
			}
		})
	}
}

func TestWrappedStreamContext(t *testing.T) {
	mockStream := &mockServerStream{
		ctx: context.Background(),
	}

	customCtx := context.WithValue(mockStream.ctx, userKey, 456)
	wrapped := &wrappedStream{
		ServerStream: mockStream,
		ctx:          customCtx,
	}

	returnedCtx := wrapped.Context()
	userID, err := GetUserIDFromContext(returnedCtx)

	assert.NoError(t, err)
	assert.Equal(t, 456, userID)
}

type mockServerStream struct {
	ctx context.Context
}

func (m *mockServerStream) SetHeader(md metadata.MD) error {
	return nil
}

func (m *mockServerStream) SendHeader(md metadata.MD) error {
	return nil
}

func (m *mockServerStream) SetTrailer(md metadata.MD) {
}

func (m *mockServerStream) Context() context.Context {
	return m.ctx
}

func (m *mockServerStream) SendMsg(v interface{}) error {
	return nil
}

func (m *mockServerStream) RecvMsg(v interface{}) error {
	return nil
}
