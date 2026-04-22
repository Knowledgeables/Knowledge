package pages

import (
	"database/sql"
)

type Language string

const (
	LanguageEN Language = "en"
	LanguageDA Language = "da"
)

type Page struct {
	Title       string       `db:"title"`
	URL         string       `db:"url"`
	Language    Language     `db:"language"`
	LastUpdated sql.NullTime `db:"last_updated" json:"lastUpdated" swaggertype:"string"`
	Content     string       `db:"content"`
}

type CrawlSignal struct {
	Query       string       `json:"query"`
	Language    Language     `json:"language"`
	SignalCount int          `json:"signalCount"`
	LastSeen    sql.NullTime `json:"lastSeen" swaggertype:"string"`
}

func NewPage(title, url, content string) Page {
	return Page{
		Title:    title,
		URL:      url,
		Language: LanguageEN,
		Content:  content,
	}
}
