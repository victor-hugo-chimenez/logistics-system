package delivery

import (
	"context"
	"fmt"
	"testing"
)

type RepositoryMock struct {
}

func (r *RepositoryMock) FindById(ctx context.Context, id int) (*Delivery, error) {
	return &Delivery{
		ID:       0,
		OrderId:  "1234",
		DriverId: "1234",
	}, nil
}

func TestService_FindById(t *testing.T) {
	deliveryService := NewDeliveryService(&RepositoryMock{})
	value, err := deliveryService.FindById(context.Background(), 0)

	if err != nil {
		t.Error(fmt.Println("Expected value but received error in FindById func"))
	}

	if value == nil {
		t.Error(fmt.Println("Expected value but received nil instead"))
	}

	if value.ID != 0 {
		t.Error(fmt.Printf("Expected ID = 0 but received ID = %d", value.ID))
	}
}
