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
