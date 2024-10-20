---
Order: 4
Title: HTTP SSL Proxy
ImgSrc: https://images.unsplash.com/photo-1489875347897-49f64b51c1f8?ixid=M3w2NjYzMTJ8MHwxfHJhbmRvbXx8fHx8fHx8fDE3MjkyNzkxOTd8&ixlib=rb-4.0.3
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

go run main.go -cert=localhost+2.pem -key=localhost+2-key.pem [-p=8003] [-ps=8004] [-url=true]
```
