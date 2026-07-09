package service

import (
	"context"
	"log/slog"

	"github.com/mykyta-kravchenko98/ShortUrl/internal/cache"
	repositories "github.com/mykyta-kravchenko98/ShortUrl/internal/db/postgres"
	"github.com/mykyta-kravchenko98/ShortUrl/internal/model"
	"github.com/mykyta-kravchenko98/ShortUrl/pkg/generator"
	hashfunction "github.com/mykyta-kravchenko98/ShortUrl/pkg/hash_function"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

var tracer = otel.Tracer("shorturl/service")

type urlService struct {
	idGenerator   generator.Snowflake
	cache         cache.LRUCache
	urlRepository repositories.URLDataService
}

// URLService - interface for working with url data
type URLService interface {
	GenerateShortURL(ctx context.Context, longURL string) (shortURL string, err error)
	GetLongURL(ctx context.Context, shortURL string) (string, error)
	Ping(ctx context.Context) error
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
func (us *urlService) GenerateShortURL(ctx context.Context, longURL string) (shortURL string, err error) {
	ctx, span := tracer.Start(ctx, "urlService.GenerateShortURL")
	defer span.End()

	existRecord, err := us.urlRepository.GetByLongURL(ctx, longURL)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		slog.ErrorContext(ctx, "lookup by long url failed", "error", err)
		return shortURL, err
	}

	if existRecord.ID > 0 {
		us.cache.Put(existRecord.ShortURL, longURL)
		slog.DebugContext(ctx, "long url already shortened", "shortURL", existRecord.ShortURL)
		return existRecord.ShortURL, nil
	}

	id, err := us.idGenerator.NextID()
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return shortURL, err
	}

	shortURL = hashfunction.DecimalToBase62(id)

	newItem := model.ShortenURLModel{
		ID:       id,
		ShortURL: shortURL,
		LongURL:  longURL,
	}

	persisted, err := us.urlRepository.Create(ctx, newItem)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		slog.ErrorContext(ctx, "failed to persist shortened url", "error", err)
		return "", err
	}

	us.cache.Put(persisted.ShortURL, persisted.LongURL)
	slog.InfoContext(ctx, "shortened new url", "shortURL", persisted.ShortURL)

	return persisted.ShortURL, nil
}

// GetLongURL - scaning chache for containing shortURL key and return longURL value
// if scaning failed it looks into db and return longURL value or exception
func (us *urlService) GetLongURL(ctx context.Context, shortURL string) (string, error) {
	ctx, span := tracer.Start(ctx, "urlService.GetLongURL")
	defer span.End()

	if longURL := us.cache.Get(shortURL); longURL != "" {
		slog.DebugContext(ctx, "cache hit", "shortURL", shortURL)
		return longURL, nil
	}

	result, err := us.urlRepository.Get(ctx, shortURL)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		slog.ErrorContext(ctx, "lookup by short url failed", "shortURL", shortURL, "error", err)
		return "", err
	}

	if result.ID == 0 {
		slog.WarnContext(ctx, "short url not found", "shortURL", shortURL)
		return "", nil
	}

	us.cache.Put(result.ShortURL, result.LongURL)

	return result.LongURL, nil
}

// Ping checks that downstream dependencies (Postgres) are reachable. Used by
// the /readyz probe.
func (us *urlService) Ping(ctx context.Context) error {
	return us.urlRepository.Ping(ctx)
}
