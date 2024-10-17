---
Order: 4
Title: HTTP SSL Proxy
Summary: Basic HTTP SSL Proxy
ImgSrc: static/article2.jpg
---

# HTTP Proxy (SSL upgrade)

## Instructions

Build upon previous [HTTP Proxy project](./http-proxy.html)
and make the proxy support ssl encryption and https connections.

## Key Features

- listen on 2 different port for http & https connections

## Usage

```shell
# brew install mkcert
# mkcert -install 
# mkcert localhost 127.0.0.1 ::1
## mkcert -uninstall

go run http-proxy-ssl/main.go -cert=localhost+2.pem -key=localhost+2-key.pem [-p=8003] [-ps=8004] [-url=true]
```
