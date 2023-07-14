package order_status

import "context"

type IRepository interface {
	FindStatusByOrderId(ctx context.Context, id int) ([]OrderStatus, error)
	UpdateOrderStatus(ctx context.Context, status *OrderStatus) error
}

type Service struct {
	repository IRepository
}

type UpdateOrderStatusCommand struct {
	status          string
	source          string
	orderId         *int
	orderExternalId *string
}

func NewOrderService(repository IRepository) *Service {
	return &Service{
		repository,
	}
}

func (s *Service) FindStatusByOrderId(ctx context.Context, id int) ([]OrderStatus, error) {
	return s.repository.FindStatusByOrderId(ctx, id)
}

func (s *Service) UpdateOrderStatus(ctx context.Context, status *OrderStatus) error {

	return s.repository.UpdateOrderStatus(ctx, status)
}
