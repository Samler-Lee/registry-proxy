package server

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"registry-proxy/internal/config"
	"registry-proxy/internal/handler"
	"registry-proxy/pkg/console"

	"github.com/labstack/echo/v4"
)

var httpServer *http.Server

func startHTTPServer() {
	engine := echo.New()
	handler.Load(engine)

	tlsConfig := &tls.Config{}
	if config.Server.EnableTLS {
		cert, err := tls.LoadX509KeyPair(config.Server.TLSCertificate, config.Server.TLSKey)
		if err != nil {
			console.Log().Error("[http] 加载TLS证书失败: %s", err)
			return
		}

		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	httpServer = &http.Server{
		Addr:    config.Server.Listen,
		Handler: engine,
	}

	console.Log().Info("[http] 正在监听: %s", httpServer.Addr)
	listenFunc := httpServer.ListenAndServe
	if config.Server.EnableTLS {
		listenFunc = func() error {
			return httpServer.ListenAndServeTLS(config.Server.TLSCertificate, config.Server.TLSKey)
		}
	}

	if err := listenFunc(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		console.Log().Error("[http] 监听 %s 失败, %s", httpServer.Addr, err)
	} else {
		console.Log().Info("[http] 服务器已关闭")
	}
}

func stopHTTPServer() {
	ctx := context.Background()

	err := httpServer.Shutdown(ctx)
	if err != nil {
		console.Log().Error("[http] 服务器关闭时出错, %s", err)
	}
}
