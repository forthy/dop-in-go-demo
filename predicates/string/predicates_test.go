package string

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotEmptyString(t *testing.T) {
	assert.False(t, NotEmpty(" "))
	assert.False(t, NotEmpty(""))
	assert.True(t, NotEmpty("hello"))
}

func TestInBetween(t *testing.T) {
	assert.True(t, InBetween(1)(5)("hello"))
	assert.False(t, InBetween(1)(5)(""))
	assert.False(t, InBetween(5)(1)("hello"))
}

func TestEmailString(t *testing.T) {
	assert.True(t, ShouldBeEmail("test@example.com"))
	assert.False(t, ShouldBeEmail("invalid-email"))
}
