package delivery

import (
	"net/http"

	threadModel "github.com/forums/app/internal/thread"
	"github.com/labstack/echo/v4"
)

type Handler struct {
}

func NewThreadHandler() threadModel.ThreadHandler {
	return &Handler{}
}

func (h *Handler) CreateThread(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}

func (h *Handler) GetDetails(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}

func (h *Handler) UpdateDetails(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}

func (h *Handler) GetPosts(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}

func (h *Handler) Vote(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}
