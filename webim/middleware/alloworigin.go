package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func AddAllowOrigin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			err := next(ctx)
			resp := ctx.Response()
			if resp.Status >= http.StatusBadRequest {
				ctx.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*")
			}
			return err
		}
	}
}
