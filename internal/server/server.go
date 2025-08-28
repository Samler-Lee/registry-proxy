package server

import (
	"context"
	"registry-proxy/internal/config"
	"registry-proxy/pkg/console"
)

var cancelFunc context.CancelFunc

func Start() {
	config.Load()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		if err := startHTTPServer(); err != nil {
			console.Log().Error("%s", err)
			cancel()
		}
	}()

	cancelFunc = cancel
	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

func Stop() {
	stopHTTPServer()
	cancelFunc()
}
