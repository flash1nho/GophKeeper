package cmd

import (
	"testing"
	"context"

	"github.com/stretchr/testify/assert"
)

func TestJwtCredentialsGetRequestMetadata(t *testing.T) {
	creds := jwtCredentials{token: "test_token_123"}

	metadata, err := creds.GetRequestMetadata(context.TODO())

	assert.NoError(t, err)
	assert.NotNil(t, metadata)
	assert.Equal(t, "test_token_123", metadata["authorization"])
}

func TestJwtCredentialsRequireTransportSecurity(t *testing.T) {
	creds := jwtCredentials{token: "test_token"}

	assert.True(t, creds.RequireTransportSecurity())
}
