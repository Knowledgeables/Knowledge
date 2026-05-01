package pages

import (
	"errors"
	"testing"
)

// compile-time check
var _ PageRepository = (*mockPageRepo)(nil)

type mockPageRepo struct {
	results   []Page
	count     int
	err       error
	lastQuery string
}

func (m *mockPageRepo) Search(q string, lang Language) ([]Page, int, error) {
	m.lastQuery = q
	return m.results, m.count, m.err
}

func (m *mockPageRepo) GetAll() ([]Page, error) {
	return nil, nil
}

func (m *mockPageRepo) FindByURL(url string) (*Page, error) {
	return nil, nil
}

func (m *mockPageRepo) RecordSignal(q string, lang Language) error {
	return nil
}

func (m *mockPageRepo) GetTopSignals(limit int) ([]CrawlSignal, error) {
	return nil, nil
}

func (m *mockPageRepo) GetExistingURLs() ([]string, error) {
	return nil, nil
}

func (m *mockPageRepo) UpsertCrawlerPages(items []CrawlerIngestItem) (int, error) {
	return 0, nil
}

// helper
func newTestService(repo *mockPageRepo) *Service {
	return NewService(repo)
}


// happy path
func TestSearch_Success(t *testing.T) {
	repo := &mockPageRepo{
		results: []Page{{Title: "Docker"}},
		count:   1,
	}

	service := newTestService(repo)

	results, count, err := service.Search("docker", LanguageEN)

	

	if err != nil {
		t.Fatal(err)
	}

	if count != 1 || len(results) != 1 {
		t.Fatal("unexpected result")
	}
}

// error path
func TestSearch_RepoError(t *testing.T) {
	repo := &mockPageRepo{
		err: errors.New("db failed"),
	}

	service := newTestService(repo)

	_, _, err := service.Search("docker", LanguageEN)

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSearch_TrimQuery(t *testing.T) {
	repo := &mockPageRepo{}
	service := newTestService(repo)

	_, _, err := service.Search("   docker   ", LanguageEN)
	if err != nil {
		t.Fatal(err)
	}

	if repo.lastQuery != "docker" {
		t.Fatalf("expected trimmed query, got '%s'", repo.lastQuery)
	}
}

func TestSearch_EmptyQuery(t *testing.T) {
	service := newTestService(&mockPageRepo{})

	results, count, err := service.Search("", LanguageEN)

	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 || count != 0 {
		t.Fatal("expected empty result")
	}
}

func TestSearch_CountReturned(t *testing.T) {
	repo := &mockPageRepo{
		results: []Page{{Title: "Docker"}},
		count:   42,
	}

	service := newTestService(repo)

	_, count, _ := service.Search("docker", LanguageEN)

	if count != 42 {
		t.Fatalf("expected 42, got %d", count)
	}
}
