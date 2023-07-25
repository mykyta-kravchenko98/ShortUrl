package repositories

import (
	"database/sql"
	"fmt"

	"github.com/mykyta-kravchenko98/ShortUrl/internal/model"
)

const (
	insertDyn = `insert into "shorten_url"("id", "short_url", "long_url") values($1, $2, $3)`
)

// Interface for access to short url db data
type URLDataService interface {
	Create(shortenURLModel model.ShortenURLModel) error
	Get(hash string) (model.ShortenURLModel, error)
	GetByLongURL(longURL string) (model.ShortenURLModel, error)
}

// return StockDataService instance
func NewCurrencySnapshotDataService(db *sql.DB) URLDataService {
	iDBSvc := &urlRepository{
		database: db,
	}
	return iDBSvc
}

// currencyRepository implements StockDataService
type urlRepository struct {
	database *sql.DB
}

// Insert record in db
func (urlRepository *urlRepository) Create(shortenURLModel model.ShortenURLModel) error {
	_, err := urlRepository.database.Exec(insertDyn, shortenURLModel.ID, shortenURLModel.ShortURL, shortenURLModel.LongURL)

	return err
}

// Get record by short_url hash
func (urlRepository *urlRepository) Get(hash string) (model.ShortenURLModel, error) {
	query := fmt.Sprintf(`SELECT id, shorten_url, long_url FROM "shorten_url" WHERE short_url = '%s'`, hash)
	model := model.ShortenURLModel{}
	rows, err := urlRepository.database.Query(query)

	if err != nil {
		return model, err
	}

	defer rows.Close()
	for rows.Next() {
		var id int64
		var shortenUrl string
		var longUrl string

		err = rows.Scan(&id, &shortenUrl, &longUrl)

		if err != nil {
			return model, err
		}

		model.ID = id
		model.ShortURL = shortenUrl
		model.LongURL = longUrl
	}

	return model, nil
}

// Get record by short_url hash
func (urlRepository *urlRepository) GetByLongURL(longURL string) (model.ShortenURLModel, error) {
	query := fmt.Sprintf(`SELECT id, shorten_url, long_url FROM "shorten_url" WHERE long_url = '%s'`, longURL)
	model := model.ShortenURLModel{}
	rows, err := urlRepository.database.Query(query)

	if err != nil {
		return model, err
	}

	defer rows.Close()
	for rows.Next() {
		var id int64
		var shortenUrl string
		var longUrl string

		err = rows.Scan(&id, &shortenUrl, &longUrl)

		if err != nil {
			return model, err
		}

		model.ID = id
		model.ShortURL = shortenUrl
		model.LongURL = longUrl
	}

	return model, nil
}
