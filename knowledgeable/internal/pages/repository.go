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
	SELECT title, url, language, last_updated, content
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
		if err := rows.Scan(&p.Title, &p.URL, &p.Language, &p.LastUpdated, &p.Content); err != nil {
			return nil, 0, err
		}
		pages = append(pages, p)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return pages, count, nil
}

func (r *Repository) GetExistingURLs() ([]string, error) {
	rows, err := r.db.Query(`
		SELECT DISTINCT url
		FROM pages
		ORDER BY url
	`)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("rows close error: %v", err)
		}
	}()

	var urls []string
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return urls, nil
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

func (r *Repository) RecordSignal(query string, lang Language) error {
	normalizedQuery := strings.ToLower(strings.TrimSpace(query))
	if normalizedQuery == "" {
		return nil
	}

	_, err := r.db.Exec(`
		INSERT INTO crawl_signals (query, language, signal_count, first_seen, last_seen)
		VALUES ($1, $2, 1, NOW(), NOW())
		ON CONFLICT (query, language)
		DO UPDATE SET
			signal_count = crawl_signals.signal_count + 1,
			last_seen = NOW()
	`, normalizedQuery, string(lang))

	return err
}

func (r *Repository) GetTopSignals(limit int) ([]CrawlSignal, error) {
	rows, err := r.db.Query(`
		SELECT query, language, signal_count, last_seen
		FROM crawl_signals
		ORDER BY signal_count DESC, last_seen DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("rows close error: %v", err)
		}
	}()

	var signals []CrawlSignal
	for rows.Next() {
		var s CrawlSignal
		if err := rows.Scan(&s.Query, &s.Language, &s.SignalCount, &s.LastSeen); err != nil {
			return nil, err
		}
		signals = append(signals, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return signals, nil
}

func (r *Repository) UpsertCrawlerPages(items []CrawlerIngestItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	stmt, err := tx.Prepare(`
		INSERT INTO pages (title, url, language, content, last_updated)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (url)
		DO UPDATE SET
			title = EXCLUDED.title,
			language = EXCLUDED.language,
			content = EXCLUDED.content,
			last_updated = NOW()
	`)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("stmt close error: %v", err)
		}
	}()

	upserted := 0
	for _, item := range items {
		if _, execErr := stmt.Exec(item.Title, item.URL, item.Language, item.Content); execErr != nil {
			_ = tx.Rollback()
			return 0, execErr
		}
		upserted++
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return upserted, nil
}
