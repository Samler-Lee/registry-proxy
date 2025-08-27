FROM golang:1.25-alpine AS builder

ADD . /app
WORKDIR /app

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

RUN go build -o registry-proxy .

FROM alpine:latest

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories

ENV TZ=Asia/Shanghai
RUN apk update \
	&& apk add tzdata \
	&& echo "${TZ}" > /etc/timezone \
	&& ln -sf /usr/share/zoneinfo/${TZ} /etc/localtime \
	&& rm /var/cache/apk/*

COPY --from=builder /app/registry-proxy /app/registry-proxy
WORKDIR /app

EXPOSE 8000
ENTRYPOINT ["./registry-proxy", "serve"]