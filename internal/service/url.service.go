package service

import (
	"github.com/mykyta-kravchenko98/ShortUrl/internal/cache"
	repositories "github.com/mykyta-kravchenko98/ShortUrl/internal/db/postgres"
	"github.com/mykyta-kravchenko98/ShortUrl/internal/model"
	"github.com/mykyta-kravchenko98/ShortUrl/pkg/generator"
	hashfunction "github.com/mykyta-kravchenko98/ShortUrl/pkg/hash_function"
)

type urlService struct {
	idGenerator   generator.Snowflake
	cache         cache.LRUCache
	urlRepository repositories.URLDataService
}

type UrlService interface {
	GenerateShortUrl(longUrl string) (shortUrl string, err error)
	GetLongUrl(shortUrl string) (string, error)
}

func NewUrlService(idGenerator generator.Snowflake, cache cache.LRUCache, rep repositories.URLDataService) UrlService {
	return &urlService{
		idGenerator:   idGenerator,
		cache:         cache,
		urlRepository: rep,
	}
}

func (us *urlService) GenerateShortUrl(longUrl string) (shortUrl string, err error) {
	existRecord, err := us.urlRepository.GetByLongURL(longUrl)

	if err != nil {
		return shortUrl, err
	}

	if existRecord.Id > 0 {
		us.cache.Put(shortUrl, longUrl)

		return existRecord.ShortURL, err
	}

	id, err := us.idGenerator.NextID()

	if err != nil {
		return shortUrl, err
	}

	shortUrl = hashfunction.DecimalToBase62(id)

	newItem := model.ShortenURLModel{
		Id:       id,
		ShortURL: shortUrl,
		LongURL:  longUrl,
	}

	err = us.urlRepository.Create(newItem)

	if err != nil {
		return "", err
	}

	us.cache.Put(shortUrl, longUrl)

	return shortUrl, err
}

func (us *urlService) GetLongUrl(shortUrl string) (string, error) {
	longUrl := us.cache.Get(shortUrl)

	if longUrl != "" {
		return longUrl, nil
	}

	result, err := us.urlRepository.Get(shortUrl)

	if err != nil {
		return longUrl, err
	}

	us.cache.Put(result.ShortURL, result.LongURL)

	return result.LongURL, err
}
