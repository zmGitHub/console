package middleware

import (
	"github.com/labstack/echo/v4"

	"bitbucket.org/forfd/custm-chat/webim/common/logging"
)

type RefreshAgentOnlineStatus func(userToken, entID string) error

func RefreshOnline(refreshFn RefreshAgentOnlineStatus) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			authStr, err := extractor(ctx)
			if err != nil {
				return err
			}

			entID := ctx.Get(AgentEntIDKey).(string)
			if authStr != "" {
				if err := refreshFn(authStr, entID); err != nil {
					log.Logger.Warnf("RefreshAgentOnlineStatus error: %v", err)
				}
			}

			return next(ctx)
		}
	}
}
