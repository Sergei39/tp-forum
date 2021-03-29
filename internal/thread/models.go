package thread

import (
	"github.com/labstack/echo/v4"
)

type ThreadHandler interface {
	CreateThread(c echo.Context) error
	GetDetails(c echo.Context) error
	UpdateDetails(c echo.Context) error
	GetPosts(c echo.Context) error
	Vote(c echo.Context) error
}
