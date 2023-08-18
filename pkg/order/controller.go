package order

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	order_item "logistics_system/pkg/order/item"
	order_status "logistics_system/pkg/order/status"
	"net/http"
	"strconv"
)

type BaseService interface {
	FindById(ctx context.Context, id int) (*Order, error)
	FindAll(ctx context.Context) ([]Order, error)
	UpdateOrder(ctx context.Context, order *Order) (*Order, error)
	CreateOrder(ctx context.Context, order *Order) error
	DeleteById(ctx context.Context, id int) error
}

type ItemService interface {
	FindItemByOrderId(ctx context.Context, id int) ([]order_item.OrderItem, error)
	CreateOrderItem(ctx context.Context, item *order_item.OrderItem) error
}

type StatusService interface {
	FindStatusByOrderId(ctx context.Context, id int) (*order_status.OrderStatus, error)
	UpdateOrderStatusHistory(ctx context.Context, status *order_status.OrderStatus) error
}

type Controller struct {
	orderService       BaseService
	orderItemService   ItemService
	orderStatusService StatusService
}

func NewController(orderService BaseService, orderItemService ItemService, orderStatusService StatusService) *Controller {
	return &Controller{
		orderService,
		orderItemService,
		orderStatusService,
	}
}

func (c *Controller) FindById(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idParam)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Error parsing order ID")
		return
	}

	order, err := c.orderService.FindById(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, fmt.Sprintf("Error getting order by id: %d", id))
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"statusCode": 200,
		"result":     order,
	})

}

func (c *Controller) FindItemByOrderId(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")

	if idParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Order Id not provided")
		return
	}

	id, err := strconv.Atoi(idParam)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Error parsing order ID")
		return
	}

	order, err := c.orderItemService.FindItemByOrderId(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, fmt.Sprintf("Error getting order by id: %d", id))
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"statusCode": 200,
		"result":     order,
	})
}

func (c *Controller) FindStatusByOrderId(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")

	if idParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Order Id not provided")
		return
	}

	id, err := strconv.Atoi(idParam)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Error parsing order ID")
		return
	}

	order, err := c.orderStatusService.FindStatusByOrderId(r.Context(), id)
	fmt.Println(err)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, fmt.Sprintf("Error getting order by id: %d", id))
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"statusCode": 200,
		"result":     order,
	})
}

func (c *Controller) FindAll(w http.ResponseWriter, r *http.Request) {
	deliveries, err := c.orderService.FindAll(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, "Error getting orders")
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"statusCode": 200,
		"result":     deliveries,
	})

}

func (c *Controller) CreateOrder(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("could not read body: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Could not read body")
		return
	}

	var order Order
	if err := json.Unmarshal(body, &order); err != nil {
		fmt.Printf("could not read body: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Could not read body")
		return
	}

	err = c.orderService.CreateOrder(r.Context(), &order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, fmt.Sprintf("Could not create order: %s\n", err))
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"statusCode": 200,
	})

}

func (c *Controller) UpdateOrder(w http.ResponseWriter, r *http.Request) {

	var order *Order
	err := json.NewDecoder(r.Body).Decode(&order)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Error parsing order")
		return
	}

	updatedOrder, err := c.orderService.UpdateOrder(r.Context(), order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, fmt.Sprintf("Error updating order by id: %d", order.ID))
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"statusCode": 200,
		"result":     updatedOrder,
	})
}

func (c *Controller) CreateOrderItem(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("could not read body: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Could not read body")
		return
	}

	var orderItem order_item.OrderItem
	if err := json.Unmarshal(body, &orderItem); err != nil {
		fmt.Printf("could not read body: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Could not read body")
		return
	}

	err = c.orderItemService.CreateOrderItem(r.Context(), &orderItem)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, fmt.Sprintf("Could not create order: %s\n", err))
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"statusCode": 200,
	})

}

func (c *Controller) CreateOrderStatus(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("could not read body: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Could not read body")
		return
	}

	var orderStatus order_status.OrderStatus
	if err := json.Unmarshal(body, &orderStatus); err != nil {
		fmt.Printf("could not read body: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Could not read body")
		return
	}

	err = c.orderStatusService.UpdateOrderStatusHistory(r.Context(), &orderStatus)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, fmt.Sprintf("Could not create order: %s\n", err))
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"statusCode": 200,
	})

}

func (c *Controller) DeleteById(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idParam)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Error parsing driver ID")
		return
	}

	if err := c.orderService.DeleteById(r.Context(), id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, fmt.Sprintf("Error getting driver by id: %d", id))
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"statusCode": 200,
	})
}

func (c *Controller) HandleOrderRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if r.URL.Query().Has("id") {
			c.FindById(w, r)
			return
		}
		c.FindAll(w, r)
		return

	case http.MethodPost:
		c.CreateOrder(w, r)
		return

	case http.MethodPut:
		c.UpdateOrder(w, r)

	case http.MethodDelete:
		c.DeleteById(w, r)
		return

	default:
		http.Error(w, "Método não suportado", http.StatusMethodNotAllowed)
	}
}

func (c *Controller) HandleOrderItemRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		c.FindItemByOrderId(w, r)
		return

	case http.MethodPost:
		c.CreateOrderItem(w, r)
		return

	default:
		http.Error(w, "Método não suportado", http.StatusMethodNotAllowed)
	}
}

func (c *Controller) HandleOrderStatusRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		c.FindStatusByOrderId(w, r)
		return

	case http.MethodPost:
		c.CreateOrderStatus(w, r)
		return

	default:
		http.Error(w, "Método não suportado", http.StatusMethodNotAllowed)
	}
}
