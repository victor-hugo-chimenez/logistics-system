package delivery

import "context"

type DeliveryRepository interface {
	FindById(ctx context.Context, id int) (*Delivery, error)
}

type Service struct {
	repository DeliveryRepository
}

func NewDeliveryService(repository DeliveryRepository) *Service {
	return &Service{
		repository,
	}
}

func (s *Service) FindById(ctx context.Context, id int) (*Delivery, error) {
	return s.repository.FindById(ctx, id)
}
