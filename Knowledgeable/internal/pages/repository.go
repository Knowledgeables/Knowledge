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
		SELECT title, url, language
		FROM pages
		WHERE language = ?
		AND LOWER(title) LIKE LOWER(?)
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
		)
		if err != nil {
			return nil, err
		}

		pages = append(pages, p)
	}

	return pages, nil
}
func (r *Repository) FindByURL(url string) (*Page, error) {

	row := r.db.QueryRow(`
		SELECT title, url, language, content
		FROM pages
		WHERE url = ?
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



