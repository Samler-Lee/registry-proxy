[中文文档](https://github.com/Samler-Lee/registry-proxy/blob/master/README.md)

# Registry Proxy

A lightweight container registry proxy service

## ✨ Features

- Support for multiple upstream
- Support for private registry login authentication
- No disk caching, friendly to devices with limited disk capacity
- Low memory consumption
- Support for Let's Encrypt automatic certificate management

## ⚙️ Configuration

A `config.toml` file is provided in the project root directory, which serves as the default configuration file for this project. You can also obtain it by directly running `registry-proxy serve`.

### proxy.binding

This configuration item is a mapping between domain names and upstream addresses. When accessing with the corresponding domain name, the request will be proxied to the corresponding upstream address.

Taking the mapping relationship in the default configuration as an example, if you access `docker.registry-proxy.localhost`, the service will forward your request to `registry-1.docker.io`.

### server.listen

The address and port that the HTTP service listens on

Default value: `:80`

### server.logLevel

The log level for console output. You can adjust it to `DEBUG` to view some details of request forwarding, but it will generate a large amount of logs.

Default value: `INFO`

### server.tls.enable

Whether to enable TLS

Default value: `false`

### server.tls.listen

The address and port that the HTTPS service listens on

Default value: `:443`

### server.tls.useLetsEncrypt

Whether to enable automatic certificate management based on `Let's Encrypt`. This configuration is mutually exclusive with `server.tls.certPath` and `server.tls.keyPath`. If this configuration is enabled, the latter will be disabled.

**Note: If this configuration is enabled, please ensure that ports `80` and `443` can correctly access this service, otherwise Let's Encrypt will refuse to issue certificates. Reference: [Challenge Types](https://letsencrypt.org/docs/challenge-types/).**

Default value: `true`

### server.tls.certPath

TLS certificate file path

Default value: `server.crt`

### server.tls.keyPath

TLS certificate key file path

Default value: `server.key`

## 🛠️ Deployment

This project was developed and tested using Go 1.25 during development. It is recommended to use Go 1.25 or above for compilation and deployment.

This project supports the following deployment methods. You can choose according to your preference:

### Docker

Execute the following command in the project root directory to build the container image
```shell
docker build -t registry-proxy:latest .
```

Start the container
```shell
docker run -itd -p 8000:8000 -v ./config.toml:/app/config.toml --restart=always registry-proxy:latest
```

### Docker Compose (Recommended)

A `docker-compose.yml` file is provided in the project root directory, which you can use for Docker Compose deployment.

First, refer to the Docker section to build the container image, then execute the following command to start the service
```shell
docker-compose up -d
```

### Binary

Compile the binary file
```shell
go build -o registry-proxy
```

Start the server
```shell
./registry-proxy serve
```

## 📖 Usage Recommendations

During the design and testing of this project, the deployment architecture used was `nginx` ---> `docker compose` ---> `registry-proxy`. It is recommended that you also adopt a similar architecture for deployment, which makes it convenient for you to reuse ports 80 and 443 and manage certificates uniformly, unless you plan to deploy only this service on the device.

## ⚖️ License

This project is distributed under GPLv3 and uses the following open source components during development:
- [echo](https://github.com/labstack/echo)
- [viper](https://github.com/spf13/viper)
- [cobra](https://github.com/spf13/cobra)
- [color](https://github.com/fatih/color)
