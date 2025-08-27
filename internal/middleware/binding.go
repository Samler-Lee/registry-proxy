package middleware

import (
	"net/http"
	"registry-proxy/internal/config"
	"registry-proxy/pkg/console"
	"strings"

	"github.com/labstack/echo/v4"
)

func DomainBinding(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		hostname := strings.Split(c.Request().Host, ":")[0]
		console.Log().Debug("current hostname: %s", hostname)

		val, exists := config.Proxy.Binding[hostname]
		if !exists {
			return echo.NewHTTPError(http.StatusNotFound, "endpoint binding not found")
		}

		c.Set("endpoint", val)
		return next(c)
	}
}
