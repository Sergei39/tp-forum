package middleware

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/forums/app/models"
	"github.com/forums/utils/logger"
	"github.com/labstack/echo/v4"
)

func LogMiddleware(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {
		start := time.Now()
		requestID := fmt.Sprintf("%016x", rand.Int())[:10]
		c.Set("request_id", requestID)
		ctx := models.GetContext(c)

		logger.Middleware().Info(ctx, logger.Fields{
			"url":           c.Request().URL,
			"method":        c.Request().Method,
			"remote_addr":   c.Request().RemoteAddr,
			"server_status": c.Response().Status,
		})

		result := next(c)

		logger.Middleware().Info(ctx, logger.Fields{
			"work_time": time.Since(start),
		})

		return result
	}
}
