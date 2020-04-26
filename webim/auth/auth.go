package auth

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

func ParseUserToken(secret, tokenStr string) (entID, agentID string, err error) {
	claims := &Claims{}
	token, jwtErr := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
		}
		return []byte(secret), nil
	})

	if jwtErr == nil && token.Valid {
		claims := token.Claims.(*Claims)
		entID, agentID = claims.EntID, claims.UserID
		return
	}

	return "", "", &echo.HTTPError{
		Code:     http.StatusUnauthorized,
		Message:  "invalid or expired jwt",
		Internal: err,
	}
}
