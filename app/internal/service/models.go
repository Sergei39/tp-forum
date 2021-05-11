package service

import (
	"context"

	"github.com/forums/app/models"
	"github.com/forums/utils/response"
	"github.com/labstack/echo/v4"
)

type ServiceHandler interface {
	ClearDb(c echo.Context) error
	StatusDb(c echo.Context) error
}

type ServiceUsecase interface {
	ClearDb(ctx context.Context) error
	StatusDb(ctx context.Context) (response.Response, error)
}

type ServiceRepo interface {
	ClearDb(ctx context.Context) error
	StatusDb(ctx context.Context) (*models.InfoStatus, error)
}
