package helpers

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestGetStructKeys(t *testing.T) {
	type TestStruct struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Age   int    `json:"age"`
	}

	keys := GetStructKeys(TestStruct{})

	assert.NotEmpty(t, keys)
	assert.Len(t, keys, 3)

	fieldNames := make(map[string]bool)
	for _, key := range keys {
		fieldNames[key.Key] = true
	}

	assert.True(t, fieldNames["name"])
	assert.True(t, fieldNames["email"])
	assert.True(t, fieldNames["age"])
}

func TestGetStructKeysWithIgnored(t *testing.T) {
	type TestStruct struct {
		Name   string `json:"name"`
		Secret string `json:"secret" ignore:"true"`
	}

	keys := GetStructKeys(TestStruct{})

	fieldNames := make(map[string]bool)
	for _, key := range keys {
		fieldNames[key.Key] = true
	}

	assert.True(t, fieldNames["name"])
	assert.False(t, fieldNames["secret"])
}

func TestGetStructKeysWithPointer(t *testing.T) {
	type TestStruct struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	keys := GetStructKeys(&TestStruct{})

	assert.NotEmpty(t, keys)
	assert.Len(t, keys, 2)
}

func TestArgsParse(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String("login", "testuser", "Login")
	cmd.Flags().String("password", "testpass", "Password")
	_ = cmd.Flags().Parse([]string{"--login", "user123", "--password", "pass123"})

	id, data, _, err := ArgsParse(cmd)

	assert.NoError(t, err)
	assert.Equal(t, 0, id)
	assert.NotNil(t, data)
}

func TestArgsParseWithID(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().Int("id", 0, "ID")
	cmd.Flags().String("content", "test", "Content")

	_ = cmd.Flags().Parse([]string{"--id", "42", "--content", "test data"})

	id, data, _, err := ArgsParse(cmd)

	assert.NoError(t, err)
	assert.Equal(t, 42, id)
	assert.NotNil(t, data)
}
