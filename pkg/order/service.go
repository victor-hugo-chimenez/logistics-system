package order

import "context"

type BaseRepository interface {
	FindById(ctx context.Context, id int) (*Order, error)
	CreateOrder(ctx context.Context, order *Order) error
	FindAll(ctx context.Context) ([]Order, error)
	UpdateOrder(ctx context.Context, order *Order) (*Order, error)
	DeleteById(ctx context.Context, id int) error
}

type Service struct {
	repository BaseRepository
}

func NewOrderService(repository BaseRepository) *Service {
	return &Service{
		repository,
	}
}

func (s *Service) FindById(ctx context.Context, id int) (*Order, error) {
	return s.repository.FindById(ctx, id)
}

func (s *Service) FindAll(ctx context.Context) ([]Order, error) {
	return s.repository.FindAll(ctx)
}

func (s *Service) CreateOrder(ctx context.Context, order *Order) error {
	return s.repository.CreateOrder(ctx, order)
}

func (s *Service) UpdateOrder(ctx context.Context, order *Order) (*Order, error) {
	return s.repository.UpdateOrder(ctx, order)
}

func (s *Service) DeleteById(ctx context.Context, id int) error {
	return s.repository.DeleteById(ctx, id)
}
