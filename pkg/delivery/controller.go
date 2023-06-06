package delivery

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type DeliveryService interface {
	FindById(ctx context.Context, id int) (*Delivery, error)
}

type Controller struct {
	service DeliveryService
}

func NewController(service DeliveryService) *Controller {
	return &Controller{
		service,
	}
}

func (c *Controller) FindById(response http.ResponseWriter, request *http.Request) {
	idParam := request.URL.Query().Get("id")

	id, err := strconv.Atoi(idParam)

	if err != nil {
		buffer := bytes.NewBufferString("Eroooou!").Bytes()
		response.Write(buffer)
		return
	}

	delivery, err := c.service.FindById(request.Context(), id)
	if err != nil {
		buffer := bytes.NewBufferString("Deu ruim!").Bytes()
		response.Write(buffer)
		return
	}

	json.NewEncoder(response).Encode(map[string]interface{}{
		"statusCode": 200,
		"result":     delivery,
	})

}

func (c Controller) FindAll(response http.ResponseWriter, request *http.Request) {

}

func (c Controller) HandleRequest(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		if request.URL.Query().Has("id") {
			c.FindById(response, request)
			return
		} else {
			c.FindAll(response, request)
			return
		}
	case http.MethodPost:
		fmt.Println("post")
	case http.MethodPut:
		fmt.Println("put")
	case http.MethodDelete:
		fmt.Println("delete")
	default:
		http.Error(response, "Método não suportado", http.StatusMethodNotAllowed)
	}
}
