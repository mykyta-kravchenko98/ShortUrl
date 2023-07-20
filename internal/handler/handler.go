package handler

import "github.com/mykyta-kravchenko98/ShortUrl/internal/service"

type Handler struct {
	urlService service.UrlService
}

func NewHandler(urlService service.UrlService) *Handler {
	return &Handler{
		urlService: urlService,
	}
}
