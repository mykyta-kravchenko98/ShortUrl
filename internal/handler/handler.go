package handler

import "github.com/mykyta-kravchenko98/ShortUrl/internal/service"

// main Handler for application
type Handler struct {
	urlService service.URLService
}

//Init method for creating new Handler
func NewHandler(urlService service.URLService) *Handler {
	return &Handler{
		urlService: urlService,
	}
}
