package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/mykyta-kravchenko98/ShortUrl/internal/model"
)

const (
	insertQuery = `
		INSERT INTO "shorten_url"("id", "short_url", "long_url")
		VALUES ($1, $2, $3)
		ON CONFLICT (long_url) DO UPDATE SET long_url = EXCLUDED.long_url
		RETURNING id, short_url, long_url`
	getByHashQuery    = `SELECT id, short_url, long_url FROM "shorten_url" WHERE short_url = $1`
	getByLongURLQuery = `SELECT id, short_url, long_url FROM "shorten_url" WHERE long_url = $1`
)

// URLDataService is a interface for access to short url db data
type URLDataService interface {
	Create(ctx context.Context, shortenURLModel model.ShortenURLModel) (model.ShortenURLModel, error)
	Get(ctx context.Context, hash string) (model.ShortenURLModel, error)
	GetByLongURL(ctx context.Context, longURL string) (model.ShortenURLModel, error)
	Ping(ctx context.Context) error
}

// NewCurrencySnapshotDataService its method for creating instance of urlRepository and return URLDataService interface
func NewCurrencySnapshotDataService(db *sql.DB) URLDataService {
	iDBSvc := &urlRepository{
		database: db,
	}
	return iDBSvc
}

// urlRepository implements URLDataService
type urlRepository struct {
	database *sql.DB
}

// Create inserts a new short url record, or - if another request already
// inserted the same long_url first - returns that existing record instead.
// Uses a parameterized query (previously this used fmt.Sprintf string
// interpolation, which was a SQL injection vector).
func (r *urlRepository) Create(ctx context.Context, shortenURLModel model.ShortenURLModel) (model.ShortenURLModel, error) {
	var m model.ShortenURLModel
	row := r.database.QueryRowContext(ctx, insertQuery, shortenURLModel.ID, shortenURLModel.ShortURL, shortenURLModel.LongURL)
	if err := row.Scan(&m.ID, &m.ShortURL, &m.LongURL); err != nil {
		return m, err
	}
	return m, nil
}

// Get looks up a record by its short_url hash. Returns a zero-value model
// (ID == 0) and a nil error when nothing is found, matching prior behavior.
func (r *urlRepository) Get(ctx context.Context, hash string) (model.ShortenURLModel, error) {
	return r.queryRow(ctx, getByHashQuery, hash)
}

// GetByLongURL looks up a record by its original long url.
func (r *urlRepository) GetByLongURL(ctx context.Context, longURL string) (model.ShortenURLModel, error) {
	return r.queryRow(ctx, getByLongURLQuery, longURL)
}

func (r *urlRepository) queryRow(ctx context.Context, query, arg string) (model.ShortenURLModel, error) {
	var m model.ShortenURLModel

	row := r.database.QueryRowContext(ctx, query, arg)
	err := row.Scan(&m.ID, &m.ShortURL, &m.LongURL)
	if errors.Is(err, sql.ErrNoRows) {
		return m, nil
	}
	if err != nil {
		return m, err
	}

	return m, nil
}

// Ping verifies the database connection is alive; used by the /readyz probe.
func (r *urlRepository) Ping(ctx context.Context) error {
	return r.database.PingContext(ctx)
}
