package pages

import "database/sql"

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
	defer rows.Close()

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

	return pages, nil
}

func (r *Repository) Search(query string, lang Language) ([]Page, error) {

	rows, err := r.db.Query(`
		SELECT title, url, language, last_updated, content
		FROM pages
		WHERE language = ?
		AND content LIKE ?
	`,
		lang,
		"%"+query+"%",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

	return pages, nil
}

