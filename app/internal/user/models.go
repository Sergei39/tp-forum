package user

import (
	"context"

	"github.com/forums/app/models"
	"github.com/forums/utils/response"
	"github.com/labstack/echo/v4"
)

type UserHandler interface {
	CreateUser(c echo.Context) error
	GetUser(c echo.Context) error
	UpdateUser(c echo.Context) error
}

type UserUsecase interface {
	CreateUser(ctx context.Context, user models.User) (response.Response, error)
	GetUserByName(ctx context.Context, name string) (response.Response, error)
	UpdateUser(ctx context.Context, user models.User) (response.Response, error)
}

type UserRepo interface {
	CreateUser(ctx context.Context, user models.User) (int, error)
	GetUserByName(ctx context.Context, name string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, user models.User) (id int, err error)
}
