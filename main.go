package main

import (
	"github.com/mykyta-kravchenko98/ShortUrl/internal/cache"
	"github.com/mykyta-kravchenko98/ShortUrl/internal/handler"
	"github.com/mykyta-kravchenko98/ShortUrl/internal/router"
	"github.com/mykyta-kravchenko98/ShortUrl/internal/service"
	"github.com/mykyta-kravchenko98/ShortUrl/pkg/generator"
)

func main() {
	r := router.New()

	v1 := r.Group("/api/v1")

	c := cache.InitLRUCache(100)
	idGen, err := generator.NewSnowflake(1, 1)
	if err != nil {
		r.Logger.Fatal(err)
	}

	urlService := service.NewUrlService(idGen, &c)

	//d := db.New()
	//db.AutoMigrate(d)

	h := handler.NewHandler(urlService)
	h.Register(v1)
	r.Logger.Fatal(r.Start("127.0.0.1:8585"))
}
