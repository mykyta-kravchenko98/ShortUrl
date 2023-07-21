package handler

import (
	"github.com/labstack/echo/v4"
)

func (h *Handler) Register(v1 *echo.Group) {
	data := v1.Group("/data")
	data.POST("/shorten", h.Shorten)

	general := v1.Group("")
	general.GET("/:hash", h.GetLongUrl)
}
