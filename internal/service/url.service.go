package service

import (
	"errors"

	"github.com/mykyta-kravchenko98/ShortUrl/internal/cache"
	"github.com/mykyta-kravchenko98/ShortUrl/pkg/generator"
	hashfunction "github.com/mykyta-kravchenko98/ShortUrl/pkg/hash_function"
)

type urlService struct {
	idGenerator generator.Snowflake
	cache       *cache.LRUCache
}

type UrlService interface {
	GenerateShortUrl(longUrl string) (shortUrl string, err error)
	GetLongUrl(shortUrl string) (longUrl string, err error)
}

func NewUrlService(idGenerator generator.Snowflake, cache *cache.LRUCache) UrlService {
	return &urlService{
		idGenerator: idGenerator,
		cache:       cache,
	}
}

func (us *urlService) GenerateShortUrl(longUrl string) (shortUrl string, err error) {
	id, err := us.idGenerator.NextID()

	if err != nil {
		return shortUrl, err
	}

	shortUrl = hashfunction.DecimalToBase62(id)

	us.cache.Put(shortUrl, longUrl)

	return shortUrl, err
}

func (us *urlService) GetLongUrl(shortUrl string) (longUrl string, err error) {
	shortUrl = us.cache.Get(shortUrl)

	if shortUrl == "" {
		return shortUrl, errors.New("Record not found")
	}

	return shortUrl, err
}
