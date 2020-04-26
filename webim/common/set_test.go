package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNonRepeatElements(t *testing.T) {
	ast := assert.New(t)
	elements := []string{"9", "9", "1", "2", "3", "3", "2", "4", "4", "6", "8", "8"}
	res := NonRepeatElements(elements)
	ast.Equal([]string{"9", "1", "2", "3", "4", "6", "8"}, res)
}
