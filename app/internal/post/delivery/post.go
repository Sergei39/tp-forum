package delivery

import (
	"net/http"
	"strconv"

	postModel "github.com/forums/app/internal/post"
	"github.com/forums/app/models"
	"github.com/forums/utils/errors"
	"github.com/forums/utils/logger"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	postUsecase postModel.PostUsecase
}

func NewPostHandler(postUsecase postModel.PostUsecase) postModel.PostHandler {
	return &Handler{
		postUsecase: postUsecase,
	}
}

func (h *Handler) GetDetails(c echo.Context) error {
	ctx := models.GetContext(c)

	related := new(models.RequestPost)
	if err := c.Bind(related); err != nil {
		sendErr := errors.New(http.StatusBadRequest, err.Error())
		logger.Delivery().Error(ctx, sendErr)
		return c.NoContent(sendErr.Code())
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	related.Id = id
	logger.Delivery().Info(ctx, logger.Fields{"request data": *related})

	response, err := h.postUsecase.GetDetails(ctx, *related)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(response.Code(), response.Body())
}

func (h *Handler) UpdateDetails(c echo.Context) error {
	ctx := models.GetContext(c)

	message := new(models.MessagePostRequest)
	if err := c.Bind(message); err != nil {
		sendErr := errors.New(http.StatusBadRequest, err.Error())
		logger.Delivery().Error(ctx, sendErr)
		return c.NoContent(sendErr.Code())
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	message.Id = id
	logger.Delivery().Info(ctx, logger.Fields{"request data": *message})

	response, err := h.postUsecase.UpdateMessage(ctx, *message)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(response.Code(), response.Body())
}
