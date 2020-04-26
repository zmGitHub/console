package middleware

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"bitbucket.org/forfd/custm-chat/webim/handler/monitor"
)

type (
	PrometheusConfig struct {
		Skipper middleware.Skipper
	}
)

var (
	DefaultPrometheusConfig = PrometheusConfig{
		Skipper: middleware.DefaultSkipper,
	}
)

func NewMetric() echo.MiddlewareFunc {
	return NewMetricWithConfig(DefaultPrometheusConfig)
}

func NewMetricWithConfig(config PrometheusConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultPrometheusConfig.Skipper
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()
			res := c.Response()
			start := time.Now()

			if err := next(c); err != nil {
				c.Error(err)
			}
			uri := req.URL.Path
			status := strconv.Itoa(res.Status)
			elapsed := time.Since(start).Seconds()
			bytesOut := float64(res.Size)

			monitor.RequestTotalCount.WithLabelValues(status, req.Method, req.Host, uri).Inc()
			monitor.RequestDurationInSec.WithLabelValues(req.Method, req.Host, uri).Observe(elapsed)
			monitor.ResponseSizeInBytes.Observe(bytesOut)
			return nil
		}
	}
}
