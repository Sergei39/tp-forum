package delivery

import (
	"net/http"

	forumModel "github.com/forums/app/internal/forum"
	"github.com/forums/app/models"
	"github.com/forums/utils/errors"
	"github.com/forums/utils/logger"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	forumUsecase forumModel.ForumUsecase
}

func NewForumHandler(usecase forumModel.ForumUsecase) forumModel.ForumHandler {
	return &Handler{
		forumUsecase: usecase,
	}
}

func (h *Handler) CreateForum(c echo.Context) error {
	ctx := models.GetContext(c)

	newForum := new(models.Forum)
	if err := c.Bind(newForum); err != nil {
		sendErr := errors.New(http.StatusBadRequest, err.Error())
		logger.Delivery().Error(ctx, sendErr)
		return c.NoContent(sendErr.Code())
	}
	logger.Delivery().Info(ctx, logger.Fields{"request data": *newForum})

	response, err := h.forumUsecase.CreateForum(ctx, *newForum)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(response.Code(), response.Body())
}

func (h *Handler) GetDetails(c echo.Context) error {
	ctx := models.GetContext(c)

	slug := c.Param("slug")
	logger.Delivery().Info(ctx, logger.Fields{"request data": slug})

	response, err := h.forumUsecase.GetForumBySlug(ctx, slug)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(response.Code(), response.Body())
}

func (h *Handler) GetUsers(c echo.Context) error {
	ctx := models.GetContext(c)

	slug := c.Param("slug")
	forumUsers := new(models.ForumUsers)
	if err := c.Bind(forumUsers); err != nil {
		sendErr := errors.New(http.StatusBadRequest, err.Error())
		logger.Delivery().Error(ctx, sendErr)
		return c.NoContent(sendErr.Code())
	}
	forumUsers.Slug = slug
	logger.Delivery().Info(ctx, logger.Fields{"request data": *forumUsers})

	response, err := h.forumUsecase.GetUsers(ctx, *forumUsers)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(response.Code(), response.Body())
}

func (h *Handler) GetThreads(c echo.Context) error {
	ctx := models.GetContext(c)

	slug := c.Param("slug")
	forumThreads := new(models.ForumThreads)
	if err := c.Bind(forumThreads); err != nil {
		sendErr := errors.New(http.StatusBadRequest, err.Error())
		logger.Delivery().Error(ctx, sendErr)
		return c.NoContent(sendErr.Code())
	}
	forumThreads.Slug = slug
	logger.Delivery().Info(ctx, logger.Fields{"request data": *forumThreads})

	response, err := h.forumUsecase.GetThreads(ctx, *forumThreads)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(response.Code(), response.Body())
}
