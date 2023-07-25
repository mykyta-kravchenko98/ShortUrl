package handler

import (
	"net/http"
	"net/url"

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
// @Router /{hash} [get]
func (h *Handler) GetLongUrl(c echo.Context) error {
	hash := c.Param("hash")

	longUrl, err := h.urlService.GetLongURL(hash)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.Redirect(http.StatusPermanentRedirect, longUrl)
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

	if !isValidURL(req.LongUrl) {
		return c.JSON(http.StatusBadRequest, "")
	}

	shortUrl, err := h.urlService.GenerateShortURL(req.LongUrl)

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

func isValidURL(inputURL string) bool {
	// Parse the URL
	u, err := url.Parse(inputURL)
	if err != nil {
		return false
	}

	// Check if the Scheme and Host are not empty
	if u.Scheme == "" || u.Host == "" {
		return false
	}

	// URL is valid
	return true
}
