package delivery

import (
	"net/http"

	serviceModel "github.com/forums/internal/service"
	"github.com/labstack/echo/v4"
)

type Handler struct {
}

func NewServiceHandler() serviceModel.ServiceHandler {
	return &Handler{}
}

func (h *Handler) ClearDb(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}

func (h *Handler) StatusDb(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}
