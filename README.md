[English Version](https://github.com/Samler-Lee/registry-proxy/blob/master/README_en.md)

# Registry Proxy

一个轻量级容器镜像代理服务

## ✨ 特性

- 支持多上游代理
- 支持私有仓库登录认证
- 支持全路径代理（代理 push 等操作）
- 无磁盘缓存，对小磁盘容量设备友好
- 较低的内存占用
- 支持 Let's Encrypt 证书自动管理

## ⚙️ 配置

项目根目录中提供了一个 `config.toml` 文件，它就是该项目的默认配置文件，你也可以直接运行 `registry-proxy serve` 来获得。

### proxy.coverAll

是否启用全路径代理，开启后，你可以使用它代理 push 操作，也可以使用其它有关的API，参考：[HTTP API V2](https://distribution.github.io/distribution/spec/api/)

**注意：如果你需要代理push，那么 `proxy.binding` 中的上游地址必须为官方仓库地址**

默认值：`false`

### proxy.binding

该配置项是域名和上游地址的关系映射，使用对应域名访问时会将请求代理到对应的上游地址。

以默认配置中的映射关系为例，如果你访问了 `docker.registry-proxy.localhost`，那么服务将会把你的请求转发至 `registry-1.docker.io`。

### server.listen

HTTP 服务监听的地址和端口

默认值：`:80`

### server.logLevel

控制台输出的日志等级，你可以调整至 `DEBUG` 查看请求转发的一些细节，但会产生大量日志。

默认值：`INFO`

### server.tls.enable

是否开启TLS

默认值：`false`

### server.tls.listen

HTTPS 服务监听的地址和端口

默认值：`:443`

### server.tls.useLetsEncrypt

是否启用基于 `Let's Encrypt` 的自动证书管理功能，该配置与 `server.tls.certPath` 和 `server.tls.keyPath` 互斥，如果启用该配置，后者将失效。

**注意：如果启用该配置，请确保 `80` 或 `443` 端口能够正确访问到本服务，否则 Let's Encrypt 会拒绝签发证书，参考：[验证方式](https://letsencrypt.org/zh-cn/docs/challenge-types/)。**

默认值：`true`

### server.tls.certPath

TLS证书文件路径

默认值：`server.crt`

### server.tls.keyPath

TLS证书密钥文件路径

默认值：`server.key`

## 🛠️ 部署

本项目在开发过程中使用 Go 1.25 进行开发和测试，建议使用 Go 1.25 及以上版本进行编译和部署

本项目支持以下部署方式，你可以根据你的喜好进行部署：

### Docker

#### 在本项目根目录下执行以下命令来构建容器镜像，或者直接使用预先构建好的公共镜像
```shell
docker build -t registry-proxy:latest .
```

#### 启动容器
```shell
docker run -itd -p 8000:80 -p 8443:443 -v ./config.toml:/app/config.toml --restart=always registry-proxy:latest
```

或

```shell
docker run -itd -p 8000:80 -p 8443:443 -v ./config.toml:/app/config.toml --restart=always ghcr.io/samler-lee/registry-proxy:latest
```

这样一来，容器将使用系统的 `8000` 和 `8443` 端口进行监听，你可以访问 `http://docker.registry-proxy.localhost:8000/v2/` 进行测试。

### Docker Compose（推荐）

在本项目根目录中提供了一个 `docker-compose.yml` 文件，你可以使用它进行 Docker Compose 方式的部署。

首先参考 Docker 章节中的构建容器镜像，或者直接使用预先构建好的公共镜像，然后执行以下命令来启动服务
```shell
docker-compose up -d
```

**公共镜像地址：`ghcr.io/samler-lee/registry-proxy:latest`，如果要使用，记得修改 `docker-compose.yml` 中的 `image` 字段。**

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

本项目在设计和测试过程中，采用的部署架构是 `nginx` ---> `docker compose` ---> `registry-proxy`，建议你也采用类似的架构进行部署，这样方便你复用80、443端口以及对证书的统一管理，除非你打算在设备上只部署该服务。

## ⚖️ 许可证

本项目采用 GPLv3 进行分发，并在开发过程中使用了以下开源组件：
- [echo](https://github.com/labstack/echo)
- [viper](https://github.com/spf13/viper)
- [cobra](https://github.com/spf13/cobra)
- [color](https://github.com/fatih/color)