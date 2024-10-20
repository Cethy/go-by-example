package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"http-proxy/handler"
	httpmiddleware "http-server-middleware/http-middleware"
	"io"
	"log"
	"net/http"
	"time"
)

var ctx = context.Background()

func GetProxyCacheHandler(urlMode bool) httpmiddleware.HandlerFunc {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return func(w http.ResponseWriter, r *http.Request) (status int, err error) {
		destUrl, err := handler.GetUrl(r, urlMode)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return http.StatusBadRequest, err
		}

		val, err := rdb.Get(ctx, destUrl).Result()

		if errors.Is(err, redis.Nil) {
			res, err := http.Get(destUrl)
			if err != nil {
				http.Error(w, "Wrong parameter format or bad reply from target destination", http.StatusBadRequest)
				return http.StatusBadRequest, err
			}
			read, err := io.ReadAll(res.Body)

			if err != nil {
				panic(err)
			}
			val = string(read)
			err = rdb.Set(ctx, destUrl, val, time.Duration(1)*time.Minute).Err()
			if err != nil {
				panic(err)
			}
		} else if err != nil {
			panic(err)
		} else {
			log.Println("CACHE HIT!", destUrl)
		}

		_, err = fmt.Fprintf(w, val)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return http.StatusBadRequest, err
		}
		return http.StatusOK, nil
	}
}
