package middleware

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"bitbucket.org/forfd/custm-chat/webim/common/logging"
)

// Log logging middleware
func Log() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			req := ctx.Request()
			if req.URL.Path == "/api/v1/status/sync" || req.URL.Path == "/metrics" {
				if err := next(ctx); err != nil {
					ctx.Error(err)
				}
				return nil
			}

			start := time.Now()
			if err := next(ctx); err != nil {
				ctx.Error(err)
			}
			stop := time.Now()

			p := req.URL.Path
			if p == "" {
				p = "/"
			}

			fields := map[string]interface{}{
				"time_rfc3339":  time.Now().Format(time.RFC3339),
				"remote_ip":     ctx.RealIP(),
				"host":          req.Host,
				"uri":           req.RequestURI,
				"method":        req.Method,
				"path":          p,
				"latency":       strconv.FormatInt(int64(stop.Sub(start)), 10),
				"latency_human": stop.Sub(start).String(),
			}

			log.Logger.WithFields(fields).Info()
			return nil
		}
	}
}
