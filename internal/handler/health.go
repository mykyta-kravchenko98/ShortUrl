package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Healthz is a liveness probe: if the process can respond at all, it's alive.
// Must never depend on downstream services (DB, cache, etc).
func (h *Handler) Healthz(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

// Readyz is a readiness probe: the pod should only receive traffic once
// dependencies (Postgres) are reachable.
func (h *Handler) Readyz(c echo.Context) error {
	ctx := c.Request().Context()

	if err := h.urlService.Ping(ctx); err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"status": "unavailable",
			"error":  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "ready"})
}
