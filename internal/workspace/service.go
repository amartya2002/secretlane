package workspace

import "errors"

type Service struct {
	repo *Repository
}

func NewService() *Service {
	return &Service{repo: NewDefaultRepository()}
}

func (s *Service) Create(name, description string, userID int) (int, error) {
	// Prevent duplicate workspace names per user
	count, err := s.repo.CountByNameForUser(name, userID)
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, errors.New("workspace with this name already exists")
	}

	return s.repo.CreateWorkspace(name, description, userID)
}

func (s *Service) ListForUser(userID int) ([]Workspace, error) {
	return s.repo.ListForUser(userID)
}

func (s *Service) Update(id int, name, description string, userID int) error {
	return s.repo.Update(id, name, description, userID)
}

func (s *Service) Delete(id int, userID int) error {
	return s.repo.Delete(id, userID)
}
