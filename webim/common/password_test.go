package common

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenHashedPassword(t *testing.T) {
	pwd := "1234abcd"
	v, err := GenHashedPassword([]byte(pwd))
	assert.Nil(t, err)
	log.Println(string(v), len(string(v)))

	pwd = "1234abcdef"
	v, err = GenHashedPassword([]byte(pwd))
	assert.Nil(t, err)
	log.Println(string(v), len(string(v)))

	pwd = "1234abcd@#"
	v, err = GenHashedPassword([]byte(pwd))
	assert.Nil(t, err)
	log.Println(string(v), len(string(v)))

	pwd = "@#$6789%^&ABcDeJ"
	v, err = GenHashedPassword([]byte(pwd))
	assert.Nil(t, err)
	log.Println(string(v), len(string(v)))
}

func TestValidatePassword(t *testing.T) {
	pwd := "1234abcd"
	v, err := GenHashedPassword([]byte(pwd))
	assert.Nil(t, err)
	hashedPwd := v

	assert.True(t, ValidatePassword(hashedPwd, []byte(pwd)))

	pwd = "1w2e3rasdf@#$"
	hashedPwd, err = GenHashedPassword([]byte(pwd))
	assert.Nil(t, err)
	assert.True(t, ValidatePassword(hashedPwd, []byte(pwd)))
}
