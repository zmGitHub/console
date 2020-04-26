package common

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// AddAllowOriginHeader add cors headers
var AddAllowOriginHeader = func(resp *echo.Response) {
	resp.Header().Set(echo.HeaderAccessControlAllowOrigin, "*")
	resp.Header().Set(echo.HeaderAccessControlAllowCredentials, "true")
	methods := []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete}
	resp.Header().Set(echo.HeaderAccessControlRequestMethod, strings.Join(methods, ", "))
}
