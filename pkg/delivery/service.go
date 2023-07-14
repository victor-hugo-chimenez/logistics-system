package delivery

import "context"

type BaseRepository interface {
	FindById(ctx context.Context, id int) (*Delivery, error)
	CreateDelivery(ctx context.Context, delivery *Delivery) error
	FindAll(ctx context.Context) ([]Delivery, error)
	UpdateDelivery(ctx context.Context, delivery *Delivery) (*Delivery, error)
	DeleteById(ctx context.Context, id int) error
}

type Service struct {
	repository BaseRepository
}

func NewDeliveryService(repository BaseRepository) *Service {
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

func (s *Service) UpdateDelivery(ctx context.Context, delivery *Delivery) (*Delivery, error) {
	return s.repository.UpdateDelivery(ctx, delivery)
}

func (s *Service) DeleteById(ctx context.Context, id int) error {
	return s.repository.DeleteById(ctx, id)
}
