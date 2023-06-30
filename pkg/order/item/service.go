package order_item

import "context"

type IRepository interface {
	FindItemByOrderId(ctx context.Context, id int) ([]OrderItem, error)
	CreateOrderItem(ctx context.Context, item *OrderItem) error
}

type Service struct {
	repository IRepository
}

func NewOrderService(repository IRepository) *Service {
	return &Service{
		repository,
	}
}

func (s *Service) FindItemByOrderId(ctx context.Context, id int) ([]OrderItem, error) {
	return s.repository.FindItemByOrderId(ctx, id)
}

func (s *Service) CreateOrderItem(ctx context.Context, item *OrderItem) error {
	return s.repository.CreateOrderItem(ctx, item)
}
