package common

import (
	"golang.org/x/crypto/bcrypt"
)

// GenHashedPassword 生成加密的用户密码
func GenHashedPassword(userPwd []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(userPwd, bcrypt.DefaultCost)
}

// ValidatePassword 验证用户密码
func ValidatePassword(hashedPwd, userPwd []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPwd, userPwd)
	return err == nil
}
