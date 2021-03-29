package delivery

import (
	"net/http"

	userModel "github.com/forums/internal/user"
	"github.com/labstack/echo/v4"
)

type Handler struct {
}

func NewUserHandler() userModel.UserHandler {
	return &Handler{}
}

func (h *Handler) CreateUser(c echo.Context) error {
	username := c.Param("username")
	return c.JSON(http.StatusOK, username)
}

func (h *Handler) GetUser(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}

func (h *Handler) UpdateUser(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}
