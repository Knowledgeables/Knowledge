//go:build integration

package pages

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/knowledge?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}

	for i := 0; i < 10; i++ {
		if err := db.Ping(); err == nil {
			break
		}
		time.Sleep(1 * time.Second)
		if i == 9 {
			t.Fatalf("db not ready")
		}
	}

	_, err = db.Exec(`TRUNCATE pages RESTART IDENTITY`)
	if err != nil {
		t.Fatalf("truncate: %v", err)
	}

	t.Cleanup(func() {
		_ = db.Close()
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

	handler := NewHandler(
		NewService(NewRepository(db)),
		func() *template.Template { return nil },
	)

	req := httptest.NewRequest(http.MethodGet, "/api/search?q=Go+Integration+Testing&language=en", nil)
	rr := httptest.NewRecorder()

	handler.SearchAPI(rr, req)

	res := rr.Result()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}

	var resp SearchResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(resp.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(resp.Results))
	}

	if resp.Results[0].Title != "Go Integration Testing" {
		t.Errorf("expected title 'Go Integration Testing', got %q", resp.Results[0].Title)
	}
}