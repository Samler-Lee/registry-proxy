package server

import "registry-proxy/internal/config"

func Start() {
	config.Load()

	startHTTPServer()
}

func Stop() {
	stopHTTPServer()
}
