package driver

import "context"

type IRepository interface {
	FindById(ctx context.Context, id int) (*Driver, error)
	FindAll(ctx context.Context) ([]Driver, error)
	CreateDriver(ctx context.Context, driver *Driver) error
}

type Service struct {
	repository IRepository
}

func NewDriverService(repository IRepository) *Service {
	return &Service{
		repository,
	}
}

func (s *Service) FindById(ctx context.Context, id int) (*Driver, error) {
	return s.repository.FindById(ctx, id)
}

func (s *Service) FindAll(ctx context.Context) ([]Driver, error) {
	return s.repository.FindAll(ctx)
}

func (s *Service) CreateDriver(ctx context.Context, driver *Driver) error {
	return s.repository.CreateDriver(ctx, driver)
}

func (s *Service) UpdateById(ctx context.Context, id int) (*Driver, error) {
	return nil, nil
}
