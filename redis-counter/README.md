---
Order: 5
Title: Redis counter
Summary: Redis counter
ImgSrc: static/article3.jpg
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
