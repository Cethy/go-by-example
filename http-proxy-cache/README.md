---
Order: 6
Title: HTTP Proxy Cache
Summary: Basic Proxy Cache
ImgSrc: static/article1.jpg
---

# HTTP Proxy (Cache upgrade)

## Instructions

Build upon previous [HTTP SSL Proxy project](./http-proxy-ssl.html)
and make the proxy cache the response to return it without polling the source next time.

## Key Features

- use redis as cache service

## Usage

```shell
# docker pull redis
# docker run -d --name redis-cache -p 6379:6379 redis
# go get github.com/redis/go-redis/v9

go run http-proxy-cache/main.go -cert=localhost+2.pem -key=localhost+2-key.pem [-p=8005] [-ps=8006] [-url=true]
```
