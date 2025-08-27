package config

var (
	Server = &server{
		Listen:         ":8000",
		LogLevel:       "INFO",
		EnableTLS:      false,
		TLSCertificate: "server.crt",
		TLSKey:         "server.key",
	}

	Proxy = &proxy{
		Binding: map[string]string{
			"docker.registry-proxy.localhost": "https://registry-1.docker.io",
			"ghcr.registry-proxy.localhost":   "https://ghcr.io",
			"gcr.registry-proxy.localhost":    "https://gcr.io",
			"quay.registry-proxy.localhost":   "https://quay.io",
			"k8s.registry-proxy.localhost":    "https://registry.k8s.io",
		},
	}
)
