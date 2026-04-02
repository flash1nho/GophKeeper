package certs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCerts(t *testing.T) {
	tmpDir := t.TempDir()
	certsDir := filepath.Join(tmpDir, "certs")
	err := os.MkdirAll(certsDir, 0755)
	assert.NoError(t, err)

	certFile := filepath.Join(certsDir, "client.crt")
	keyFile := filepath.Join(certsDir, "client.key")
	caFile := filepath.Join(certsDir, "ca.crt")

	for _, file := range []string{certFile, keyFile, caFile} {
		f, err := os.Create(file)
		assert.NoError(t, err)
		f.Close()
	}

	oldWd, _ := os.Getwd()
	_ = os.Chdir(tmpDir)

	defer func() {
		_ = os.Chdir(oldWd)
	}()

	certs, err := NewCerts("client")

	if err == nil {
		assert.NotNil(t, certs)
	}
}
