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
}

type PostUsecase interface {
	GetDetails(ctx context.Context, request models.RequestPost) (response.Response, error)
	UpdateMessage(ctx context.Context, request models.MessagePostRequest) (response.Response, error)
}

type PostRepo interface {
	GetPost(ctx context.Context, id int) (*models.Post, error)
	UpdateMessage(ctx context.Context, request models.MessagePostRequest) error
}
