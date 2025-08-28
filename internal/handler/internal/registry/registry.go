package registry

import (
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"registry-proxy/internal/pkg/util"
	"registry-proxy/pkg/console"
	"time"

	"github.com/labstack/echo/v4"
)

var client = &http.Client{
	Transport: &http.Transport{
		TLSHandshakeTimeout: 10 * time.Second,
		DialContext: (&net.Dialer{
			Timeout: 30 * time.Second,
		}).DialContext,
		IdleConnTimeout:       90 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 10 * time.Second,
	},
}

func LoginProxy(c echo.Context) error {
	endpoint := c.Param("endpoint")
	if endpoint == "" {
		return errors.New("endpoint is empty")
	}

	endpoint, err := url.PathUnescape(endpoint)
	if err != nil {
		return err
	}

	method := c.Request().Method
	requestURL := endpoint + "?" + c.QueryString()

	console.Log().Debug("%s: %s", method, requestURL)
	request, err := http.NewRequest(method, requestURL, c.Request().Body)
	if err != nil {
		return err
	}

	util.SetRequestHeader(c.Request(), request.Header)
	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	for key, values := range resp.Header {
		c.Response().Header()[key] = values
	}

	c.Response().WriteHeader(resp.StatusCode)
	_, err = io.Copy(c.Response().Writer, resp.Body)
	return err
}

func GetRoot(c echo.Context) error {
	endpoint := c.Get("endpoint").(string)

	method := c.Request().Method
	requestURL := util.GetRequestURL(endpoint, "v2") + "/"

	console.Log().Debug("%s: %s", method, requestURL)
	request, err := http.NewRequest(c.Request().Method, requestURL, c.Request().Body)
	if err != nil {
		return err
	}

	util.SetRequestHeader(c.Request(), request.Header)
	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode == http.StatusUnauthorized {
		util.AuthenticateRedirect(resp.Header, util.GetRealScheme(c), c.Request().Host)
	}

	for key, values := range resp.Header {
		c.Response().Header()[key] = values
	}

	c.Response().WriteHeader(resp.StatusCode)
	_, err = io.Copy(c.Response().Writer, resp.Body)
	return err
}

func GetManifests(c echo.Context) error {
	endpoint := c.Get("endpoint").(string)
	repo := c.Param("repo")
	name := c.Param("name")
	reference := c.Param("reference")

	method := c.Request().Method
	requestURL := util.GetRequestURL(endpoint, "v2", repo, name, "manifests", reference)

	console.Log().Debug("%s: %s", method, requestURL)
	request, err := http.NewRequest(c.Request().Method, requestURL, c.Request().Body)
	if err != nil {
		return err
	}

	util.SetRequestHeader(c.Request(), request.Header)
	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode == http.StatusUnauthorized {
		util.AuthenticateRedirect(resp.Header, util.GetRealScheme(c), c.Request().Host)
	}

	for key, values := range resp.Header {
		c.Response().Header()[key] = values
	}

	c.Response().WriteHeader(resp.StatusCode)
	_, err = io.Copy(c.Response().Writer, resp.Body)
	return err
}

func GetBlobs(c echo.Context) error {
	endpoint := c.Get("endpoint").(string)
	repo := c.Param("repo")
	name := c.Param("name")
	digest := c.Param("digest")

	method := c.Request().Method
	requestURL := util.GetRequestURL(endpoint, "v2", repo, name, "blobs", digest)

	console.Log().Debug("%s: %s", method, requestURL)
	request, err := http.NewRequest(method, requestURL, nil)
	if err != nil {
		return err
	}

	util.SetRequestHeader(c.Request(), request.Header)
	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode == http.StatusUnauthorized {
		util.AuthenticateRedirect(resp.Header, util.GetRealScheme(c), c.Request().Host)
	}

	for key, values := range resp.Header {
		c.Response().Header()[key] = values
	}

	c.Response().WriteHeader(resp.StatusCode)
	_, err = io.Copy(c.Response().Writer, resp.Body)
	return err
}
