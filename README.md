# Registry Proxy

一个轻量级容器镜像代理服务

## ✨ 特性

- 支持多上游代理
- 支持私有仓库登录认证
- 无磁盘缓存，对小磁盘容量设备友好
- 较低的内存占用

## ⚙️ 配置

项目根目录中提供了一个 `config.toml` 文件，它就是该项目的默认配置文件，你也可以直接运行 `registry-proxy serve` 来获得。

### proxy.binding

该配置项是域名和上游地址的关系映射，使用对应域名访问时会将请求代理到对应的上游地址。

以默认配置中的映射关系为例，如果你访问了 `docker.registry-proxy.localhost`，那么服务将会把你的请求转发至 `registry-1.docker.io`。

### server.listen

服务监听的地址和端口，默认值为 `:8000`，表示监听所有网卡的 8000 端口，如果你使用 Docker 或 Docker Compose 进行部署，请确保映射了该端口。

### server.logLevel

控制台输出的日志等级，默认值为 `INFO`，你可以调整至 `DEBUG` 查看请求转发的一些细节，但会产生大量日志。

### server.enableTLS

是否开启TLS，如果您选择开启TLS，请注意提供正确的 `server.tlsCertificate` 和 `server.tlsKey`

### server.tlsCertificate

TLS证书文件路径

### server.tlsKey

TLS证书密钥文件路径

## 🛠️ 部署

本项目在开发过程中使用 Go 1.25 进行开发和测试，建议使用 Go 1.25 及以上版本进行编译和部署

本项目支持以下部署方式，你可以根据你的喜好进行部署：

### Docker

在本项目根目录下执行以下命令来构建容器镜像
```shell
docker build -t registry-proxy:latest .
```

启动容器
```shell
docker run -itd -p 8000:8000 -v ./config.toml:/app/config.toml --restart=always registry-proxy:latest
```

### Docker Compose（推荐）

再本项目根目录中提供了一个`docker-compose.yml`，你可以使用它进行 Docker Compose 方式的部署。

首先参考 Docker 章节中的构建容器镜像，然后执行以下命令来启动服务
```shell
docker-compose up -d
```

### 二进制

编译二进制文件
```shell
go build -o registry-proxy
```

启动服务器
```shell
./registry-proxy serve
```

## 📖 使用建议

本项目在设计和测试过程中，采用的部署架构是 nginx --- docker compose，建议你也采用类似的架构进行部署，这样方便你对https证书进行管理。

## ⚖️ 许可证

本项目采用 GPLv3 进行分发，并在开发过程中使用了以下开源组件：
- [echo](https://github.com/labstack/echo)
- [viper](https://github.com/spf13/viper)
- [cobra](https://github.com/spf13/cobra)
- [color](https://github.com/fatih/color)