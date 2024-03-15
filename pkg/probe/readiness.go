package probe

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// DependencyChecker is an interface that should be implemented by any service
// that the application depends on, such as a database or message queue.
// The Check method should return an error if the service is not ready.
type DependencyChecker interface {
	Check() error
}

// RedisClientHealthChecker is a struct that holds a pointer to a Redis client,
// and implements the DependencyChecker interface.
type RedisHealthChecker struct {
	Client *redis.Client
}

// Check checks the health of the Redis connection by executing PING command.
func (r *RedisHealthChecker) Check() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.Client.Ping(ctx).Result()
	return err
}

// ReadinessProbe returns an HTTP handler that checks the readiness of the given dependencies.
// It returns a 200 status code if all dependencies are ready, and a 503 status code if any
// dependency is not ready.
func ReadinessProbe(logger *zap.Logger, dependencies ...DependencyChecker) echo.HandlerFunc {
	return func(c echo.Context) error {
		for _, d := range dependencies {
			if err := d.Check(); err != nil {
				logger.Error("dependency is not ready", zap.Error(err))
				// If a dependency is not ready, return a 503 status code.
				return c.NoContent(http.StatusServiceUnavailable)
			}
		}
		// If all dependencies are ready, return a 200 status code.
		return c.NoContent(http.StatusOK)
	}
}
