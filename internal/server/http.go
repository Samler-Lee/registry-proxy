package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"registry-proxy/internal/config"
	"registry-proxy/internal/handler"
	"registry-proxy/pkg/console"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/acme/autocert"
)

var (
	httpServer      *http.Server
	httpsServer     *http.Server
	autoCertManager *autocert.Manager
)

func listenHTTP() error {
	console.Log().Info("[http] 正在监听: %s", httpServer.Addr)
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("[http] 监听 %s 失败, %s", httpServer.Addr, err)
	}

	console.Log().Info("[http] 服务器已关闭")
	return nil
}

func listenHTTPS() error {
	if !config.Server.TLS.Enable {
		return nil
	}

	var certPath, keyPath string
	if !config.Server.TLS.UseLetsEncrypt {
		certPath = config.Server.TLS.CertPath
		keyPath = config.Server.TLS.KeyPath
	}

	console.Log().Info("[https] 正在监听: %s", httpsServer.Addr)
	if err := httpsServer.ListenAndServeTLS(certPath, keyPath); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("[https] 监听 %s 时错误, %s", httpsServer.Addr, err)
	}

	console.Log().Info("[https] 服务器已关闭")
	return nil
}

func startHTTPServer() error {
	engine := echo.New()
	handler.Load(engine)

	if config.Server.TLS.Enable {
		if config.Server.TLS.UseLetsEncrypt {
			autoCertManager = &autocert.Manager{
				Prompt: autocert.AcceptTOS,
				Cache:  autocert.DirCache(".acme-cache"),
			}

			engine.GET("/.well-known/acme-challenge/*", echo.WrapHandler(autoCertManager.HTTPHandler(nil)))
		} else if config.Server.TLS.CertPath == "" || config.Server.TLS.KeyPath == "" {
			return errors.New("[https] server.tls.certPath 和 server.tls.keyPath 需要配置")
		}

		httpsServer = &http.Server{
			Addr:    config.Server.TLS.Listen,
			Handler: engine,
		}

		if autoCertManager != nil {
			httpsServer.TLSConfig = autoCertManager.TLSConfig()
		}
	}

	httpServer = &http.Server{
		Addr:    config.Server.Listen,
		Handler: engine,
	}

	errCh := make(chan error, 2)
	go func() {
		errCh <- listenHTTP()
	}()

	go func() {
		errCh <- listenHTTPS()
	}()

	return <-errCh
}

func stopHTTPServer() {
	ctx := context.Background()

	err := httpServer.Shutdown(ctx)
	if err != nil {
		console.Log().Error("[http] 服务器关闭时出错, %s", err)
	}

	if httpsServer != nil {
		err = httpsServer.Shutdown(ctx)
		if err != nil {
			console.Log().Error("[http] 服务器关闭时出错, %s", err)
		}
	}
}
