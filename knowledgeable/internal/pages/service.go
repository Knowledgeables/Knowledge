package pages

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
