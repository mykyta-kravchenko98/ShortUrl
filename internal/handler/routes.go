package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

// Register is a method for registration all avalible endpoints in Handler instance
func (h *Handler) Register(v1 *echo.Group) {
	data := v1.Group("/data")
	data.Use(shortenRateLimiter())
	data.POST("/shorten", h.Shorten)

	general := v1.Group("")
	general.GET("/:hash", h.GetLongURL)
	general.GET("/status", h.GetStatus)
}

// RegisterHealth registers liveness/readiness probes on the root router
// (not under /api/v1, so they match standard k8s probe conventions).
func (h *Handler) RegisterHealth(e *echo.Echo) {
	e.GET("/healthz", h.Healthz)
	e.GET("/readyz", h.Readyz)
}

// shortenRateLimiter caps short-URL creation at ~20/min per client IP.
// This is an unauthenticated write endpoint - without a limit, anyone can
// mass-create links (spam, or use the service as a disposable open
// redirector). Deliberately scoped to /data/shorten only: GET /:hash
// (the redirect/read path) stays unlimited.
func shortenRateLimiter() echo.MiddlewareFunc {
	return middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(middleware.RateLimiterMemoryStoreConfig{
			Rate:      rate.Limit(20.0 / 60.0), // sustained ~20/min
			Burst:     20,                      // let a client use a minute's budget in one go
			ExpiresIn: 3 * time.Minute,
		}),
		// c.RealIP() trusts X-Forwarded-For/X-Real-IP by default, which is
		// only safe once this sits behind a proxy/ingress that sets those
		// headers itself and strips client-supplied ones. If this ever runs
		// with the internet reaching it directly (no LB in front), switch
		// to echo.ExtractIPDirect() instead - otherwise the limit is
		// trivially bypassed by spoofing the header.
		IdentifierExtractor: func(c echo.Context) (string, error) {
			return c.RealIP(), nil
		},
		ErrorHandler: func(c echo.Context, err error) error {
			slog.ErrorContext(c.Request().Context(), "rate limiter store error", "error", err)
			return c.JSON(http.StatusInternalServerError, nil)
		},
		DenyHandler: func(c echo.Context, identifier string, err error) error {
			slog.WarnContext(c.Request().Context(), "rate limit exceeded", "identifier", identifier)
			return c.JSON(http.StatusTooManyRequests, map[string]string{
				"error": "rate limit exceeded, try again later",
			})
		},
	})
}
