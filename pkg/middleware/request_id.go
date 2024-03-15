package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

// RequestID is a middleware that generates a unique request ID for each HTTP request.
func RequestID(header string) echo.MiddlewareFunc {
	return echoMiddleware.RequestIDWithConfig(echoMiddleware.RequestIDConfig{
		Generator: func() string {
			return uuid.New().String()
		},
		TargetHeader: header,
	})
}
