package models

import (
	"context"

	"github.com/labstack/echo/v4"
)

func GetContext(c echo.Context) context.Context {
	ctx := c.Request().Context()

	return context.WithValue(ctx, "request_id", c.Get("request_id"))
}
