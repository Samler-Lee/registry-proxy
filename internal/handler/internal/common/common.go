package common

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func EmptyResponse(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{})
}
