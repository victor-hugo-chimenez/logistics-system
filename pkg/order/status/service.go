package order_status

import "context"

type BaseRepository interface {
	FindStatusByOrderId(ctx context.Context, id int) ([]OrderStatus, error)
	UpdateOrderStatusCheckpoint(ctx context.Context, status *OrderStatus) error
	UpdateOrderStatusHistory(ctx context.Context, status *OrderStatus) error
}

type Service struct {
	repository BaseRepository
}

type UpdateOrderStatusCommand struct {
	status          string
	source          string
	orderId         *int
	orderExternalId *string
}

func NewOrderService(repository BaseRepository) *Service {
	return &Service{
		repository,
	}
}

func (s *Service) FindStatusByOrderId(ctx context.Context, id int) ([]OrderStatus, error) {
	return s.repository.FindStatusByOrderId(ctx, id)
}

func (s *Service) UpdateOrderStatusHistory(ctx context.Context, status *OrderStatus) error {

	return s.repository.UpdateOrderStatusHistory(ctx, status)
}
