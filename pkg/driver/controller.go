package driver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type IService interface {
	FindById(ctx context.Context, id int) (*Driver, error)
	FindAll(ctx context.Context) ([]Driver, error)
	UpdateById(ctx context.Context, id int) (*Driver, error)
	CreateDriver(ctx context.Context, driver *Driver) error
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
		_, _ = io.WriteString(w, "Error parsing driver ID")
		return
	}

	driver, err := c.service.FindById(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, fmt.Sprintf("Error getting driver by id: %d", id))
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"statusCode": 200,
		"result":     driver,
	})

}

func (c *Controller) FindAll(w http.ResponseWriter, r *http.Request) {
	drivers, err := c.service.FindAll(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, "Error getting drivers")
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"statusCode": 200,
		"result":     drivers,
	})
}

func (c *Controller) CreateDriver(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("could not read body: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Could not read body")
		return
	}

	var driver Driver
	if err := json.Unmarshal(body, &driver); err != nil {
		fmt.Printf("could not read body: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Could not read body")
		return
	}

	err = c.service.CreateDriver(r.Context(), &driver)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, fmt.Sprintf("Could not create driver: %s\n", err))
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
		_, _ = io.WriteString(w, "Error parsing driver ID")
		return
	}

	driver, err := c.service.UpdateById(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, fmt.Sprintf("Error getting driver by id: %d", id))
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"statusCode": 200,
		"result":     driver,
	})
}

func (c Controller) HandleDriverRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if r.URL.Query().Has("id") {
			c.FindById(w, r)
			return
		} else {
			c.FindAll(w, r)
			return
		}
	case http.MethodPost:
		c.CreateDriver(w, r)
		return
	case http.MethodPut:
		fmt.Println("put")
	case http.MethodDelete:
		fmt.Println("delete")
	default:
		http.Error(w, "Método não suportado", http.StatusMethodNotAllowed)
	}
}
