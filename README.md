# Learning Go

https://cethy.github.io/go-by-example/

## Scripts

### hello-world

```shell
go run hello-world/main.go
```

### http-server

```shell
go run http-server/main.go [-p=8001]
```

### http-server (w/ middlewares)

```shell
go run http-server-middleware/main.go [-p=8002]
```

### http-proxy

```shell
go run http-proxy/main.go [-p=8003] [-url=true]
```

### http-proxy-ssl

```shell
# brew install mkcert
# mkcert -install 
# mkcert localhost 127.0.0.1 ::1
## mkcert -uninstall

go run http-proxy-ssl/main.go -cert=localhost+2.pem -key=localhost+2-key.pem [-p=8003] [-ps=8004] [-url=true]
```

### redis-counter

```shell
# docker pull redis
# docker run -d --name redis-counter -p 6379:6379 redis
# go get github.com/redis/go-redis/v9

go run redis-counter/main.go
```

### http-proxy-cache

```shell
# docker pull redis
# docker run -d --name redis-cache -p 6379:6379 redis
# go get github.com/redis/go-redis/v9

go run http-proxy-cache/main.go -cert=localhost+2.pem -key=localhost+2-key.pem [-p=8005] [-ps=8006] [-url=true]
```

### regex

```shell
go run regex/main.go -m mf{ze}fz{e}f{}foo
```

### markdown2tailwindcss

```shell
# go get github.com/yuin/goldmark

go run markdown2tailwindcss/main.go
```


### static-website-builder

```shell
# go get github.com/fsnotify/fsnotify
# go get github.com/yuin/goldmark-meta

# build
go run static-website-generator/generator.go
# serve
go run static-website-generator/serve.go
```

### bubbletea stuff

```shell
# go get github.com/charmbracelet/bubbletea
# go get github.com/charmbracelet/bubbles/help
# go get github.com/charmbracelet/bubbles/textinput

go run cli-todolist/main.go
```

## Sources

- https://go.dev/tour/
- https://gowebexamples.com
- https://markphelps.me/posts/handling-errors-in-your-http-handlers/
- https://medium.com/@matryer/the-http-handler-wrapper-technique-in-golang-updated-bc7fbcffa702#.e4k81jxd3

## Useful stuff

### Init module

```shell
go mod init go-by-example
```

### Installing package

```shell
go get -u github.com/gorilla/mux
```

## Ongoing questions

- what modules actually means ?
- how much concurrent connection can http-server or http-proxy can handle ?
- monitoring ?
- what is a `go` routine ?