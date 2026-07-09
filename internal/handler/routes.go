package handler

import (
	"github.com/labstack/echo/v4"
)

// Register is a method for registration all avalible endpoints in Handler instance
func (h *Handler) Register(v1 *echo.Group) {
	data := v1.Group("/data")
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
