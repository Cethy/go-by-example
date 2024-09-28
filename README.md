# Learning Go

## Scripts

### hello-world

```
go run hello-world/main.go
```

### http-server

```
go run http-server/main.go [-p=8001]
```

### http-server (w/ middlewares)

```
go run http-server-middleware/main.go [-p=8002]
```

### http-proxy

```
go run http-proxy/main.go [-p=8003] [-url=true]
```

### http-proxy-ssl

```
go run http-proxy/main.go [-p=8004] [-url=true]
```

- ...

## Sources

- https://go.dev/tour/
- https://gowebexamples.com

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
