package handler

import (
	"github.com/dgrijalva/jwt-go"

	"bitbucket.org/forfd/custm-chat/webim/auth"
	"bitbucket.org/forfd/custm-chat/webim/conf"
)

func newJwtToken(claims *auth.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(conf.IMConf.JWTKeyBytes)
}

func newConnJwtToken(claims *jwt.StandardClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(conf.IMConf.JWTKeyBytes)
}
