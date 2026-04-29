package pages

import "strings"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAllPages() ([]Page, error) {
	return s.repo.GetAll()
}

func (s *Service) Search(q string, lang Language) ([]Page, int, error) {

	if q == "" {
		return []Page{}, 0, nil
	}

	return s.repo.Search(q, lang)
}

func (s *Service) FindByURL(url string) (*Page, error) {
	return s.repo.FindByURL(url)
}

func (s *Service) RecordSignal(q string, lang Language) error {
	return s.repo.RecordSignal(q, lang)
}

func (s *Service) GetTopSignals(limit int) ([]CrawlSignal, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	return s.repo.GetTopSignals(limit)
}

func (s *Service) GetExistingURLs() ([]string, error) {
	return s.repo.GetExistingURLs()
}

func (s *Service) IngestCrawlerPages(items []CrawlerIngestItem) (int, error) {
	validated := make([]CrawlerIngestItem, 0, len(items))
	for _, item := range items {
		title := strings.TrimSpace(item.Title)
		url := strings.TrimSpace(item.URL)
		content := strings.TrimSpace(item.Content)

		if title == "" || url == "" || content == "" {
			continue
		}

		language := item.Language
		if language == "" {
			language = LanguageEN
		}

		validated = append(validated, CrawlerIngestItem{
			Title:    title,
			URL:      url,
			Language: language,
			Content:  content,
		})
	}

	if len(validated) == 0 {
		return 0, nil
	}

	return s.repo.UpsertCrawlerPages(validated)
}



