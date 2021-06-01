package delivery

import (
	"net/http"

	threadModel "github.com/forums/app/internal/thread"
	"github.com/forums/app/models"
	"github.com/forums/utils/errors"
	"github.com/forums/utils/logger"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	threadUsecase threadModel.ThreadUsecase
}

func NewThreadHandler(usecase threadModel.ThreadUsecase) threadModel.ThreadHandler {
	return &Handler{
		threadUsecase: usecase,
	}
}

func (h *Handler) CreateThread(c echo.Context) error {
	ctx := models.GetContext(c)

	slug := c.Param("slug")
	newThread := new(models.Thread)
	if err := c.Bind(newThread); err != nil {
		sendErr := errors.New(http.StatusBadRequest, err.Error())
		logger.Delivery().Error(ctx, sendErr)
		return c.NoContent(sendErr.Code())
	}
	logger.Delivery().Info(ctx, logger.Fields{"request data": *newThread})

	response, err := h.threadUsecase.CreateThread(ctx, *newThread, slug)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(response.Code(), response.Body())
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