package model

// ShortenURLModel is a main domain model that representing db table structure
type ShortenURLModel struct {
	ID       int64
	ShortURL string
	LongURL  string
}
