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

// URLService - interface for working with url data
type URLService interface {
	GenerateShortURL(longURL string) (shortURL string, err error)
	GetLongURL(shortURL string) (string, error)
}

// NewURLService - create new urlService instance and returning UrlService interface for interact with it
func NewURLService(idGenerator generator.Snowflake, cache cache.LRUCache, rep repositories.URLDataService) URLService {
	return &urlService{
		idGenerator:   idGenerator,
		cache:         cache,
		urlRepository: rep,
	}
}

// GenerateShortURL - checking if Url is already saved in DB. If yes, that returns the hash
// if no, then generate Id and use it for getting hash.
// Save data in db and add into cache. After that return saved hash.
func (us *urlService) GenerateShortURL(longURL string) (shortURL string, err error) {
	existRecord, err := us.urlRepository.GetByLongURL(longURL)

	if err != nil {
		return shortURL, err
	}

	if existRecord.ID > 0 {
		us.cache.Put(shortURL, longURL)

		return existRecord.ShortURL, err
	}

	id, err := us.idGenerator.NextID()

	if err != nil {
		return shortURL, err
	}

	shortURL = hashfunction.DecimalToBase62(id)

	newItem := model.ShortenURLModel{
		ID:       id,
		ShortURL: shortURL,
		LongURL:  longURL,
	}

	err = us.urlRepository.Create(newItem)

	if err != nil {
		return "", err
	}

	us.cache.Put(shortURL, longURL)

	return shortURL, err
}

// GetLongURL - scaning chache for containing shortURL key and return longURL value
// if scaning failed it looks into db and return longURL value or exception
func (us *urlService) GetLongURL(shortURL string) (string, error) {
	longURL := us.cache.Get(shortURL)

	if longURL != "" {
		return longURL, nil
	}

	result, err := us.urlRepository.Get(shortURL)

	if err != nil {
		return longURL, err
	}

	us.cache.Put(result.ShortURL, result.LongURL)

	return result.LongURL, err
}
