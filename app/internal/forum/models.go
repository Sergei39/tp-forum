package forum

import (
	"context"

	"github.com/forums/app/models"
	"github.com/forums/utils/response"
	"github.com/labstack/echo/v4"
)

type ForumHandler interface {
	CreateForum(c echo.Context) error
	GetDetails(c echo.Context) error
	CreateThread(c echo.Context) error
	GetUsers(c echo.Context) error
	GetThreads(c echo.Context) error
}

type ForumUsecase interface {
	CreateForum(ctx context.Context, forum models.Forum) (response.Response, error)
	GetForumBySlug(ctx context.Context, slug string) (response.Response, error)
}

type UserRepo interface {
	CreateForum(ctx context.Context, forum models.Forum) (int, error)
	GetForumBySlug(ctx context.Context, title string) (*models.Forum, error)
}
