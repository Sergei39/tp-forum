package thread

import (
	"context"

	"github.com/forums/app/models"
	"github.com/forums/utils/response"

	"github.com/labstack/echo/v4"
)

type ThreadHandler interface {
	CreateThread(c echo.Context) error
	GetDetails(c echo.Context) error
	UpdateDetails(c echo.Context) error
	GetPosts(c echo.Context) error
	Vote(c echo.Context) error
}

type ThreadUsecase interface {
	CreateThread(ctx context.Context, thread models.Thread, slug string) (response.Response, error)
	// CreateUser(ctx context.Context, user models.User) (response.Response, error)
	// GetUserByName(ctx context.Context, name string) (response.Response, error)
}

type ThreadRepo interface {
	CreateThread(ctx context.Context, thread models.Thread) (int, error)
	// CreateUser(ctx context.Context, user models.User) (int, error)
	// GetUserByName(ctx context.Context, name string) (*models.User, error)
}
