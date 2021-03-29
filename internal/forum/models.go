package forum

import (
	"github.com/labstack/echo/v4"
)

type ForumHandler interface {
	CreateForum(c echo.Context) error
	GetDetails(c echo.Context) error
	CreateThread(c echo.Context) error
	GetUsers(c echo.Context) error
	GetThreads(c echo.Context) error
}
