package delivery

import (
	"net/http"

	postModel "github.com/forums/internal/post"
	"github.com/labstack/echo/v4"
)

type Handler struct {
}

func NewPostHandler() postModel.PostHandler {
	return &Handler{}
}

func (h *Handler) GetDetails(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}

func (h *Handler) UpdateDetails(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}
