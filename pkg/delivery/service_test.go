package delivery

import "testing"

type RepositoryMock struct {
}

func TestDoSomething(t *testing.T) {
	deliveryService := NewDeliveryService(&RepositoryMock{})
	value := deliveryService.DoSomething()
	if value != nil {
		t.Error("Eroou!")
	}
}
