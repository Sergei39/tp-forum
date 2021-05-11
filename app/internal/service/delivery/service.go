package delivery

import (
	"net/http"

	serviceModel "github.com/forums/app/internal/service"
	"github.com/forums/app/models"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	serviceUsecase serviceModel.ServiceUsecase
}

func NewServiceHandler(serviceUsecase serviceModel.ServiceUsecase) serviceModel.ServiceHandler {
	return &Handler{
		serviceUsecase: serviceUsecase,
	}
}

func (h *Handler) ClearDb(c echo.Context) error {
	ctx := models.GetContext(c)

	err := h.serviceUsecase.ClearDb(ctx)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) StatusDb(c echo.Context) error {
	ctx := models.GetContext(c)

	response, err := h.serviceUsecase.StatusDb(ctx)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(response.Code(), response.Body())
}
