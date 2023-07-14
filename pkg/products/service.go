package products

import "context"

type BaseRepository interface {
	FindById(ctx context.Context, id int) (*Product, error)
	CreateProduct(ctx context.Context, product *Product) error
	FindAll(ctx context.Context) ([]Product, error)
	UpdateProduct(ctx context.Context, product *Product) (*Product, error)
	DeleteById(ctx context.Context, id int) error
}

type Service struct {
	repository BaseRepository
}

func NewProductService(repository BaseRepository) *Service {
	return &Service{
		repository,
	}
}

func (s *Service) FindById(ctx context.Context, id int) (*Product, error) {
	return s.repository.FindById(ctx, id)
}

func (s *Service) FindAll(ctx context.Context) ([]Product, error) {
	return s.repository.FindAll(ctx)
}

func (s *Service) CreateProduct(ctx context.Context, product *Product) error {
	return s.repository.CreateProduct(ctx, product)
}

func (s *Service) UpdateProduct(ctx context.Context, product *Product) (*Product, error) {
	return s.repository.UpdateProduct(ctx, product)
}

func (s *Service) DeleteById(ctx context.Context, id int) error {
	return s.repository.DeleteById(ctx, id)
}
