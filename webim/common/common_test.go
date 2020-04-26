package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenPlaceHolders(t *testing.T) {
	s := GenPlaceHolders(3)
	assert.Equal(t, "?,?,?", s)

	s = GenPlaceHolders(1)
	assert.Equal(t, "?", s)
}
