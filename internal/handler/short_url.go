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
func (h *Handler) GetLongURL(c echo.Context) error {
	hash := c.Param("hash")

	longURL, err := h.urlService.GetLongURL(hash)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.Redirect(http.StatusPermanentRedirect, longURL)
}

// Shorten godoc
// @Summary Create shortURL
// @Description Create shortURL and return it
// @ID shorten
// @Tags article
// @Accept  json
// @Produce  json
// @Param shortenRequest body, longURL is required
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

	if !isValidURL(req.LongURL) {
		return c.JSON(http.StatusBadRequest, "")
	}

	shortURL, err := h.urlService.GenerateShortURL(req.LongURL)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, shortenResponse{shortURL})
}

type shortenRequest struct {
	LongURL string `json:"longURL" validate:"required"`
}

type shortenResponse struct {
	ShortURL string `json:"shortURL" validate:"required"`
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
