package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

// RequestTimeout is a middleware that returns a 503 Service Unavailable error if the request takes longer than the specified duration.
// It uses the provided duration to set the timeout.
func RequestTimeout(timeoutSeconds uint16) echo.MiddlewareFunc {
	return echoMiddleware.TimeoutWithConfig(echoMiddleware.TimeoutConfig{
		Timeout: time.Duration(timeoutSeconds) * time.Second,
	})
}
