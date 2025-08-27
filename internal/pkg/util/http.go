package util

import (
	"net/http"
	"net/url"
	"registry-proxy/pkg/console"
	"strings"

	"github.com/labstack/echo/v4"
)

func GetRealScheme(c echo.Context) string {
	if scheme := c.Request().Header.Get("X-Forwarded-Proto"); scheme != "" {
		return scheme
	}

	if scheme := c.Request().Header.Get("X-Forwarded-Scheme"); scheme != "" {
		return scheme
	}

	if c.Request().Header.Get("X-Forwarded-Ssl") == "on" {
		return "https"
	}

	if scheme := c.Request().Header.Get("X-Url-Scheme"); scheme != "" {
		return scheme
	}

	return c.Scheme()
}

func GetRequestURL(endpoint string, paths ...string) string {
	var builder strings.Builder
	builder.WriteString(endpoint)

	for _, path := range paths {
		if path != "" {
			builder.WriteByte('/')
			builder.WriteString(path)
		}
	}

	return builder.String()
}

func AuthenticateRedirect(header http.Header, scheme string, host string) {
	console.Log().Debug("scheme: %s, host: %s", scheme, host)
	if authenticate := header.Get("WWW-Authenticate"); authenticate != "" {
		if idx := strings.Index(authenticate, "realm=\""); idx != -1 {
			start := idx + 7
			if end := strings.Index(authenticate[start:], "\""); end != -1 {
				realm := authenticate[start : start+end]
				newRealm := scheme + "://" + host + "/token/proxy/" + url.QueryEscape(realm)
				newAuthenticate := authenticate[:start] + newRealm + authenticate[start+end:]
				header.Set("WWW-Authenticate", newAuthenticate)
			}
		}
	}
}

func SetRequestHeader(request *http.Request, header http.Header) {
	for key, values := range request.Header {
		header.Set(key, strings.Join(values, ","))
	}
}
