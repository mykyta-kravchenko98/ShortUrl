package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetLongUrl godoc
// @Summary Get an long url
// @Description Get an long url if it exist
// @ID get-long-url
// @Tags url
// @Accept  json
// @Produce  json
// @Param url param for extract long url
// @Success 200 {object} singleArticleResponse
// @Failure 400 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Router /{url} [get]
func (h *Handler) GetLongUrl(c echo.Context) error {
	url := c.Param("url")

	longUrl, err := h.urlService.GetLongUrl(url)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, longUrl)
}

// Shorten godoc
// @Summary Create shortUrl
// @Description Create shortUrl and return it
// @ID shorten
// @Tags article
// @Accept  json
// @Produce  json
// @Param shortenRequest body, longUrl is required
// @Success 201 {object} singleArticleResponse
// @Failure 401 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Router /shorten [post]
func (h *Handler) Shorten(c echo.Context) error {
	req := &shortenRequest{}

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, "")
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, "")
	}

	shortUrl, err := h.urlService.GenerateShortUrl(req.LongUrl)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, shortenResponse{shortUrl})
}

type shortenRequest struct {
	LongUrl string `json:"longUrl" validate:"required"`
}

type shortenResponse struct {
	ShortUrl string `json:"shortUrl" validate:"required"`
}
