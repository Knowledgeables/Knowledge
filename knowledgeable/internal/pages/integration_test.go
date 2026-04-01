package pages

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "modernc.org/sqlite"
)

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open in-memory db: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE pages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			url TEXT NOT NULL UNIQUE,
			language TEXT NOT NULL CHECK(language IN ('en', 'da')) DEFAULT 'en',
			content TEXT NOT NULL,
			last_updated DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatalf("create table: %v", err)
	}

	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Logf("failed to close db: %v", err)
		}
	})

	return db
}

func TestSearchAPI_ReturnsStubPage(t *testing.T) {
	db := newTestDB(t)

	_, err := db.Exec(`
		INSERT INTO pages (title, url, language, content)
		VALUES ('Go Integration Testing', 'https://example.com/go-testing', 'en', 'some content')
	`)
	if err != nil {
		t.Fatalf("insert stub page: %v", err)
	}

	handler := NewHandler(NewService(NewRepository(db)))

	req := httptest.NewRequest(http.MethodGet, "/api/search?q=Go+Integration+Testing&language=en", nil)
	rr := httptest.NewRecorder()

	handler.SearchAPI(rr, req)

	res := rr.Result()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}

	var results []Page
	if err := json.NewDecoder(res.Body).Decode(&results); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	if results[0].Title != "Go Integration Testing" {
		t.Errorf("expected title 'Go Integration Testing', got %q", results[0].Title)
	}
}
