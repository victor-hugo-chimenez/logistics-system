package delivery

import "context"

type IRepository interface {
	FindById(ctx context.Context, id int) (*Delivery, error)
	CreateDelivery(ctx context.Context, delivery *Delivery) error
	FindAll(ctx context.Context) ([]Delivery, error)
}

type Service struct {
	repository IRepository
}

func NewDeliveryService(repository IRepository) *Service {
	return &Service{
		repository,
	}
}

func (s *Service) FindById(ctx context.Context, id int) (*Delivery, error) {
	return s.repository.FindById(ctx, id)
}

func (s *Service) FindAll(ctx context.Context) ([]Delivery, error) {
	return s.repository.FindAll(ctx)
}

func (s *Service) CreateDelivery(ctx context.Context, delivery *Delivery) error {
	return s.repository.CreateDelivery(ctx, delivery)
}

func (s *Service) UpdateById(ctx context.Context, id int) (*Delivery, error) {
	return nil, nil
}
