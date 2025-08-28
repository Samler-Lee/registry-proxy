package config

var (
	Server = &server{
		Listen:   ":80",
		LogLevel: "INFO",
		TLS: &tls{
			Enable:         false,
			Listen:         ":443",
			UseLetsEncrypt: true,
			CertPath:       "server.crt",
			KeyPath:        "server.key",
		},
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
