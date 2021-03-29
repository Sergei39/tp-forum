package service

import (
	"github.com/labstack/echo/v4"
)

type ServiceHandler interface {
	ClearDb(c echo.Context) error
	StatusDb(c echo.Context) error
}
