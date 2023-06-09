package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type BaseService interface {
	FindById(ctx context.Context, id int) (*Delivery, error)
	FindAll(ctx context.Context) ([]Delivery, error)
	UpdateDelivery(ctx context.Context, delivery *Delivery) (*Delivery, error)
	CreateDelivery(ctx context.Context, delivery *Delivery) error
	DeleteById(ctx context.Context, id int) error
}

type Controller struct {
	service BaseService
}

func NewController(service BaseService) *Controller {
	return &Controller{
		service,
	}
}

func (c *Controller) FindById(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idParam)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Error parsing delivery ID")
		return
	}

	delivery, err := c.service.FindById(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, fmt.Sprintf("Error getting delivery by id: %d", id))
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"statusCode": 200,
		"result":     delivery,
	})

}

func (c *Controller) FindAll(w http.ResponseWriter, r *http.Request) {
	deliveries, err := c.service.FindAll(r.Context())
	fmt.Println(err)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, "Error getting deliveries")
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"statusCode": 200,
		"result":     deliveries,
	})

}

func (c *Controller) CreateDelivery(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("could not read body: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Could not read body")
		return
	}

	var delivery Delivery
	if err := json.Unmarshal(body, &delivery); err != nil {
		fmt.Printf("could not read body: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Could not read body")
		return
	}

	err = c.service.CreateDelivery(r.Context(), &delivery)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, fmt.Sprintf("Could not create delivery: %s\n", err))
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"statusCode": 200,
	})

}

func (c *Controller) UpdateDelivery(w http.ResponseWriter, r *http.Request) {

	var delivery *Delivery
	err := json.NewDecoder(r.Body).Decode(&delivery)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Error parsing delivery")
		return
	}

	updatedDelivery, err := c.service.UpdateDelivery(r.Context(), delivery)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, fmt.Sprintf("Error updating delivery by id: %d", delivery.ID))
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"statusCode": 200,
		"result":     updatedDelivery,
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

	if err := c.service.DeleteById(r.Context(), id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, fmt.Sprintf("Error getting driver by id: %d", id))
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"statusCode": 200,
	})
}

func (c *Controller) HandleDeliveryRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if r.URL.Query().Has("id") {
			c.FindById(w, r)
			return
		}
		c.FindAll(w, r)
		return

	case http.MethodPost:
		c.CreateDelivery(w, r)
		return

	case http.MethodPut:
		c.UpdateDelivery(w, r)

	case http.MethodDelete:
		c.DeleteById(w, r)
		return

	default:
		http.Error(w, "Método não suportado", http.StatusMethodNotAllowed)
	}
}
