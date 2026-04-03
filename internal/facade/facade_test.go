package facade

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFacadeGetUserIDFromContext(t *testing.T) {
	f := NewFacade()

	assert.NotNil(t, f)
}

func TestNewFacade(t *testing.T) {
	f := NewFacade()
	assert.NotNil(t, f)
	assert.IsType(t, &Facade{}, f)
}
