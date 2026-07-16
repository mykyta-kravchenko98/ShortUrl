package observability

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func HTTPLoggerMiddleware() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogRemoteIP:      true,
		LogHost:          true,
		LogMethod:        true,
		LogURI:           true,
		LogRoutePath:     true,
		LogUserAgent:     true,
		LogStatus:        true,
		LogError:         true,
		LogLatency:       true,
		LogContentLength: true,
		LogResponseSize:  true,

		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			slog.InfoContext(
				c.Request().Context(),
				"HTTP request completed",
				"remote_ip", v.RemoteIP,
				"host", v.Host,
				"method", v.Method,
				"uri", v.URI,
				"route", v.RoutePath,
				"user_agent", v.UserAgent,
				"status", v.Status,
				"error", v.Error,
				"latency", v.Latency.Nanoseconds(),
				"latency_human", v.Latency.String(),
				"bytes_in", v.ContentLength,
				"bytes_out", v.ResponseSize,
			)

			return nil
		},
	})
}
