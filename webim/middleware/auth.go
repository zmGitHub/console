package middleware

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"

	"bitbucket.org/forfd/custm-chat/webim/auth"
	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/common/logging"
	"bitbucket.org/forfd/custm-chat/webim/db"
)

var (
	AgentEntIDKey = "EntID"
	AgentTokenKey = "AgentToken"
	AgentIDKey    = "AgentID"
)

var unauthorizedErr = func(resp *echo.Response, err error) *echo.HTTPError {
	common.AddAllowOriginHeader(resp)
	return &echo.HTTPError{
		Code:     http.StatusUnauthorized,
		Message:  "invalid or expired jwt",
		Internal: err,
	}
}

func AgentAuth(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		fn := func(ctx echo.Context) error {
			authStr, err := extractor(ctx)
			if err != nil {
				return err
			}

			claims := &auth.Claims{}
			token, jwtErr := jwt.ParseWithClaims(authStr, claims, func(t *jwt.Token) (interface{}, error) {
				if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
					return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
				}
				return []byte(secret), nil
			})

			if jwtErr == nil && token.Valid {
				exists, err := auth.TokenExists(db.RedisClient, token.Raw)
				if err != nil {
					return unauthorizedErr(ctx.Response(), err)
				}

				if !exists {
					return unauthorizedErr(ctx.Response(), err)
				}

				// Store user information from token into context.
				claims := token.Claims.(*auth.Claims)
				ctx.Set(AgentEntIDKey, claims.EntID)
				ctx.Set(AgentTokenKey, authStr)
				ctx.Set(AgentIDKey, claims.UserID)
				return next(ctx)
			}

			if err := auth.RemoveToken(db.RedisClient, authStr); err != nil {
				log.Logger.Warnf("remove invalid agent token err: %v", err)
			}

			return unauthorizedErr(ctx.Response(), jwtErr)
		}

		return fn
	}
}

func extractor(ctx echo.Context) (auth string, err *echo.HTTPError) {
	auth = ctx.Request().Header.Get(echo.HeaderAuthorization)
	if auth != "" {
		return
	}

	return "", unauthorizedErr(ctx.Response(), fmt.Errorf("%s", "missing or malformed jwt"))
}
