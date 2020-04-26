package middleware

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func SuperAdminAuth(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authStr, err := extractor(c)
			if err != nil {
				return err
			}

			if authStr == secret {
				return next(c)
			}

			return unauthorizedErr(c.Response(), fmt.Errorf("you are not super admin"))
		}
	}
}
