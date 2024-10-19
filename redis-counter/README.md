---
Order: 5
Title: Redis counter
ImgSrc: https://images.unsplash.com/photo-1517512006864-7edc3b933137?ixid=M3w2NjYzMTJ8MHwxfHJhbmRvbXx8fHx8fHx8fDE3MjkyNzkxOTh8&ixlib=rb-4.0.3
---

# Redis counter

## Instructions

Make a program capable of storing a counter in redis 
and increment each time its invoked.

## Usage

```shell
# docker pull redis
# docker run -d --name redis-counter -p 6379:6379 redis
# go get github.com/redis/go-redis/v9

go run redis-counter/main.go
```
