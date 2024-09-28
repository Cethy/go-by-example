package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
)

var ctx = context.Background()

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	val, err := rdb.Get(ctx, "count").Result()
	if errors.Is(err, redis.Nil) {
		val = "0"
	} else if err != nil {
		panic(err)
	}
	fmt.Println("count", val)

	// string to int
	i, err := strconv.Atoi(val)
	if err != nil {
		// ... handle error
		panic(err)
	}

	err = rdb.Set(ctx, "count", i+1, 0).Err()
	if err != nil {
		panic(err)
	}
}
