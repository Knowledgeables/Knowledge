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
	LastUpdated sql.NullTime `db:"last_updated"`
	Content     string       `db:"content"`
}

func NewPage(title, url, content string) Page {
	return Page{
		Title:    title,
		URL:      url,
		Language: LanguageEN,
		Content:  content,
	}
}
