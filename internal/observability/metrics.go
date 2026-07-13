package observability

import (
	"time"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func HTTPMetricsMiddleware() echo.MiddlewareFunc {
	meter := otel.Meter(ServiceName)

	requestCount, err := meter.Int64Counter(
		"http.server.request_count",
		metric.WithDescription("Number of HTTP requests handled"),
	)
	if err != nil {
		panic(err)
	}

	requestDuration, err := meter.Float64Histogram(
		"http.server.duration",
		metric.WithDescription("HTTP request duration"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		panic(err)
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			handlerErr := next(c)

			attrs := metric.WithAttributes(
				attribute.String("http.method", c.Request().Method),
				attribute.String("http.route", routeOrPath(c)),
				attribute.Int("http.status_code", c.Response().Status),
			)

			ctx := c.Request().Context()
			requestCount.Add(ctx, 1, attrs)
			requestDuration.Record(ctx, float64(time.Since(start).Milliseconds()), attrs)

			return handlerErr
		}
	}
}

// routeOrPath prefers the matched Echo route pattern (e.g. "/api/v1/:hash")
// over the raw request path, so metrics group by endpoint shape instead of
// exploding into a separate series per distinct short-URL hash.
func routeOrPath(c echo.Context) string {
	if p := c.Path(); p != "" {
		return p
	}
	return c.Request().URL.Path
}
