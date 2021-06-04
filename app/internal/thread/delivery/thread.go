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
	ctx := models.GetContext(c)

	slugOrId := c.Param("slug_or_id")
	logger.Delivery().Info(ctx, logger.Fields{"request data slug or id": slugOrId})

	response, err := h.threadUsecase.GetThread(ctx, slugOrId)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(response.Code(), response.Body())
}

func (h *Handler) UpdateDetails(c echo.Context) error {
	ctx := models.GetContext(c)

	slugOrId := c.Param("slug_or_id")
	newThread := new(models.Thread)
	if err := c.Bind(newThread); err != nil {
		sendErr := errors.New(http.StatusBadRequest, err.Error())
		logger.Delivery().Error(ctx, sendErr)
		return c.NoContent(sendErr.Code())
	}
	logger.Delivery().Info(ctx, logger.Fields{"request data": *newThread})

	response, err := h.threadUsecase.UpdateThread(ctx, *newThread, slugOrId)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(response.Code(), response.Body())
}

func (h *Handler) GetPosts(c echo.Context) error {
	ctx := models.GetContext(c)

	threadPosts := new(models.ThreadPosts)

	threadPosts.SlugOrId = c.Param("slug_or_id")
	threadPosts.Limit = c.QueryParam("limit")
	threadPosts.Since = c.QueryParam("since")
	threadPosts.Sort = c.QueryParam("sort")

	desc := c.QueryParam("desc")
	if desc == "false" || desc == "" {
		threadPosts.Desc = false
	} else {
		threadPosts.Desc = true
	}

	logger.Delivery().Info(ctx, logger.Fields{"request data": *threadPosts})

	response, err := h.threadUsecase.GetPosts(ctx, *threadPosts)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(response.Code(), response.Body())
}

func (h *Handler) Vote(c echo.Context) error {
	ctx := models.GetContext(c)

	slugOrId := c.Param("slug_or_id")
	vote := new(models.Vote)
	if err := c.Bind(vote); err != nil {
		sendErr := errors.New(http.StatusBadRequest, err.Error())
		logger.Delivery().Error(ctx, sendErr)
		return c.NoContent(sendErr.Code())
	}
	logger.Delivery().Info(ctx, logger.Fields{"request data": *vote})

	response, err := h.threadUsecase.AddVote(ctx, *vote, slugOrId)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(response.Code(), response.Body())
}
