package pages

import (
	"database/sql"
	"log"
	"strings"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAll() ([]Page, error) {

	rows, err := r.db.Query(`
		SELECT title, url, language, last_updated, content
		FROM pages
	`)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("rows close error: %v", err)
		}
	}()

	var pages []Page

	for rows.Next() {
		var p Page

		err := rows.Scan(
			&p.Title,
			&p.URL,
			&p.Language,
			&p.LastUpdated,
			&p.Content,
		)
		if err != nil {
			return nil, err
		}

		pages = append(pages, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return pages, nil
}

func (r *Repository) Search(query string, lang Language) ([]Page, int, error) {

	like := "%" + strings.ToLower(query) + "%"

	var count int
	err := r.db.QueryRow(`
        SELECT COUNT(*)
        FROM pages
        WHERE language = $1 AND LOWER(title) LIKE $2
    `, lang, like).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(`
        SELECT title, url
        FROM pages
        WHERE language = $1 AND LOWER(title) LIKE $2
    `, lang, like)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("rows close error: %v", err)
		}
	}()

	var pages []Page

	for rows.Next() {
		var p Page
		if err := rows.Scan(&p.Title, &p.URL); err != nil {
			return nil, 0, err
		}
		pages = append(pages, p)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return pages, count, nil
}

func (r *Repository) FindByURL(url string) (*Page, error) {

	row := r.db.QueryRow(`
		SELECT title, url, language, content
		FROM pages
		WHERE url = $1
	`, url)

	var p Page

	err := row.Scan(
		&p.Title,
		&p.URL,
		&p.Language,
		&p.Content,
	)

	if err != nil {
		return nil, err
	}

	return &p, nil
}
