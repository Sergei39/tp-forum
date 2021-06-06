package post

import (
	"context"

	"github.com/forums/app/models"
	"github.com/forums/utils/response"
	"github.com/labstack/echo/v4"
)

type PostHandler interface {
	GetDetails(c echo.Context) error
	UpdateDetails(c echo.Context) error
	CreatePosts(c echo.Context) error
}

type PostUsecase interface {
	GetDetails(ctx context.Context, request models.RequestPost) (response.Response, error)
	UpdateMessage(ctx context.Context, request models.MessagePostRequest) (response.Response, error)
	CreatePosts(ctx context.Context, posts []models.Post, slugOrId string) (response.Response, error)
}

type PostRepo interface {
	GetPost(ctx context.Context, id int) (*models.Post, error)
	UpdateMessage(ctx context.Context, request models.MessagePostRequest) error
	CreatePosts(ctx context.Context, posts []models.Post) ([]models.Post, error)
	CreateForumsUsers(ctx context.Context, posts []models.Post) error
}
