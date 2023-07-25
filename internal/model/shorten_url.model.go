package model

// Main domain model that present db table structure
type ShortenURLModel struct {
	ID       int64
	ShortURL string
	LongURL  string
}
