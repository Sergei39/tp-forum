package delivery

import (
	"net/http"

	forumModel "github.com/forums/internal/forum"
	"github.com/labstack/echo/v4"
)

type Handler struct {
}

func NewForumHandler() forumModel.ForumHandler {
	return &Handler{}
}

func (h *Handler) CreateForum(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}

func (h *Handler) GetDetails(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}

func (h *Handler) CreateThread(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}

func (h *Handler) GetUsers(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}

func (h *Handler) GetThreads(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}
