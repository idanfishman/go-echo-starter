package probe

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// LivenessProbe is a simple HTTP handler that always returns a 200 status code.
// Suppose to be used as a liveness probe in a Kubernetes deployment.
func LivenessProbe(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}
