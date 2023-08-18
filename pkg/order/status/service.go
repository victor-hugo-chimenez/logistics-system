package order_status

import (
	"context"
	"fmt"
)

type BaseRepository interface {
	FindStatusByOrderId(ctx context.Context, id int) (*OrderStatus, error)
	UpdateOrderStatusCheckpoint(ctx context.Context, id int) (*OrderStatus, error)
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

func (s *Service) FindStatusByOrderId(ctx context.Context, id int) (*OrderStatus, error) {
	return s.repository.FindStatusByOrderId(ctx, id)
}

func (s *Service) UpdateOrderStatusHistory(ctx context.Context, status *OrderStatus) error {
	return s.repository.UpdateOrderStatusHistory(ctx, status)
}

func (s *Service) UpdateOrderStatusCheckpoint(ctx context.Context, id int) (*OrderStatus, error) {
	orderStatus, err := s.repository.FindStatusByOrderId(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update order status checkpoint %s", err)
	}
	return s.repository.UpdateOrderStatusCheckpoint(ctx, orderStatus.OrderId)
}
