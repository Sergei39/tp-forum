package post

import (
	"github.com/labstack/echo/v4"
)

type PostHandler interface {
	GetDetails(c echo.Context) error
	UpdateDetails(c echo.Context) error
}
