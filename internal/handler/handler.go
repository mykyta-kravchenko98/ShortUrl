package handler

import "github.com/mykyta-kravchenko98/ShortUrl/internal/service"

// Handler its a main handler for application
type Handler struct {
	urlService service.URLService
}

//NewHandler is a init method for creating and returning new Handler instance
func NewHandler(urlService service.URLService) *Handler {
	return &Handler{
		urlService: urlService,
	}
}
