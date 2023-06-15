package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type IService interface {
	FindById(ctx context.Context, id int) (*Delivery, error)
	FindAll(ctx context.Context) ([]Delivery, error)
	UpdateById(ctx context.Context, id int) (*Delivery, error)
	CreateDelivery(ctx context.Context, delivery *Delivery) error
}

type Controller struct {
	service IService
}

func NewController(service IService) *Controller {
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

func (c *Controller) UpdateById(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idParam)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Error parsing delivery ID")
		return
	}

	delivery, err := c.service.UpdateById(r.Context(), id)
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
		fmt.Println("put")

	case http.MethodDelete:
		fmt.Println("delete")

	default:
		http.Error(w, "Método não suportado", http.StatusMethodNotAllowed)
	}
}

func (c *Controller) NewRouter() http.HandlerFunc {
	mux := http.NewServeMux()

	mux.HandleFunc("/", c.HandleDeliveryRequest)

	return mux.ServeHTTP
}
