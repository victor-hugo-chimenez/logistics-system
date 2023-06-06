package delivery

import (
	"bytes"
	"context"
	"net/http"
	"strconv"
	"strings"
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
	slicedPath := strings.Split(request.URL.Path, "/")
	idParam := strings.Join(slicedPath[1:], ",")

	id, err := strconv.Atoi(idParam)

	if err != nil {
		buffer := bytes.NewBufferString("Eroooou!").Bytes()
		response.Write(buffer)
		return
	}

	c.service.FindById(request.Context(), id)
}
