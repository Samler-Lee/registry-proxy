package handler

import (
	"errors"
	"net/http"
	"registry-proxy/internal/handler/internal/common"
	"registry-proxy/internal/handler/internal/registry"
	"registry-proxy/internal/middleware"

	"github.com/labstack/echo/v4"
)

func Load(e *echo.Echo) {
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		var httpError *echo.HTTPError
		if errors.As(err, &httpError) {
			_ = c.JSON(httpError.Code, map[string]any{
				"errors": []map[string]any{
					{"code": "INTERNAL_ERROR", "message": httpError.Message, "detail": err.Error()},
				},
			})
		} else {
			_ = c.JSON(http.StatusOK, map[string]any{
				"errors": []map[string]any{
					{"code": "INTERNAL_ERROR", "message": "runtime error", "detail": err.Error()},
				},
			})
		}
	}

	e.Any("/token/proxy/:endpoint", registry.LoginProxy)

	e.GET("/v2", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/v2/")
	})

	v2 := e.Group("/v2", middleware.DomainBinding)
	{
		v2.GET("/", common.EmptyResponse)

		v2.HEAD("/:name/manifests/:reference", registry.GetManifests)
		v2.HEAD("/:repo/:name/manifests/:reference", registry.GetManifests)
		v2.GET("/:name/manifests/:reference", registry.GetManifests)
		v2.GET("/:repo/:name/manifests/:reference", registry.GetManifests)

		v2.HEAD("/:name/blobs/:digest", registry.GetBlobs)
		v2.HEAD("/:repo/:name/blobs/:digest", registry.GetBlobs)
		v2.GET("/:name/blobs/:digest", registry.GetBlobs)
		v2.GET("/:repo/:name/blobs/:digest", registry.GetBlobs)
	}
}
