package products

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type BaseService interface {
	FindById(ctx context.Context, id int) (*Product, error)
	FindAll(ctx context.Context) ([]Product, error)
	UpdateProduct(ctx context.Context, product *Product) (*Product, error)
	CreateProduct(ctx context.Context, product *Product) error
	DeleteById(ctx context.Context, id int) error
}

type Controller struct {
	service BaseService
}

func NewProductController(service BaseService) *Controller {
	return &Controller{
		service,
	}
}

func (c *Controller) FindById(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idParam)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Error parsing product ID")
		return
	}

	product, err := c.service.FindById(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, fmt.Sprintf("Error getting product by id: %d", id))
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"statusCode": 200,
		"result":     product,
	})

}

func (c *Controller) FindAll(w http.ResponseWriter, r *http.Request) {
	products, err := c.service.FindAll(r.Context())
	fmt.Println(err)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, "Error getting products")
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"statusCode": 200,
		"result":     products,
	})

}

func (c *Controller) CreateProduct(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("could not read body: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Could not read body")
		return
	}

	var product Product
	if err := json.Unmarshal(body, &product); err != nil {
		fmt.Printf("could not read body: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Could not read body")
		return
	}

	err = c.service.CreateProduct(r.Context(), &product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, fmt.Sprintf("Could not create product: %s\n", err))
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"statusCode": 200,
	})

}

func (c *Controller) UpdateProduct(w http.ResponseWriter, r *http.Request) {

	var product *Product
	err := json.NewDecoder(r.Body).Decode(&product)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Error parsing product")
		return
	}

	updatedProduct, err := c.service.UpdateProduct(r.Context(), product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, fmt.Sprintf("Error updating product by id: %d", product.ID))
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"statusCode": 200,
		"result":     updatedProduct,
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

func (c *Controller) HandleProductRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if r.URL.Query().Has("id") {
			c.FindById(w, r)
			return
		}
		c.FindAll(w, r)
		return

	case http.MethodPost:
		c.CreateProduct(w, r)
		return

	case http.MethodPut:
		c.UpdateProduct(w, r)

	case http.MethodDelete:
		c.DeleteById(w, r)
		return

	default:
		http.Error(w, "Método não suportado", http.StatusMethodNotAllowed)
	}
}
