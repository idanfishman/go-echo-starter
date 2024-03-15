package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

// RequestLogger is a middleware that logs each request using the provided zap.Logger.
// It logs the URI, status, error (if any), latency, remote IP, host, and method of the request.
func RequestLogger(logger *zap.Logger) echo.MiddlewareFunc {
	return echoMiddleware.RequestLoggerWithConfig(echoMiddleware.RequestLoggerConfig{
		Skipper: func(ctx echo.Context) bool {
			uri := ctx.Request().URL.String()
			status := ctx.Response().Status
			return status == http.StatusOK && (uri == "/healthz" || uri == "/readyz")
		},
		LogURI:      true,
		LogStatus:   true,
		LogError:    true,
		LogLatency:  true,
		LogRemoteIP: true,
		LogHost:     true,
		LogMethod:   true,
		LogValuesFunc: func(c echo.Context, v echoMiddleware.RequestLoggerValues) error {
			if v.Error != nil {
				logger.Error("request",
					zap.String("uri", v.URI),
					zap.Int("status", v.Status),
					zap.String("error", v.Error.Error()),
					zap.Duration("latency_ns", v.Latency),
					zap.String("remote_ip", v.RemoteIP),
					zap.String("host", v.Host),
					zap.String("method", v.Method),
				)
			} else {
				logger.Info("request",
					zap.String("uri", v.URI),
					zap.Int("status", v.Status),
					zap.Duration("latency_ns", v.Latency),
					zap.String("remote_ip", v.RemoteIP),
					zap.String("host", v.Host),
					zap.String("method", v.Method),
				)
			}
			return nil
		},
	})
}
