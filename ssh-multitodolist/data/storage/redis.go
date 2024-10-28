package storage

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"ssh-multitodolist/data"
)

type Redis struct {
	client *redis.Client
	key    string
}

var ctx = context.Background()

func NewRedisClient(addr, password string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	return rdb
}

func NewRedisStorage(client *redis.Client, roomName string) *Redis {
	return &Redis{client, roomName}
}

func (p *Redis) Init() ([]data.NamedList, error) {
	var namedLists []data.NamedList

	jsonData, err := p.client.Get(ctx, p.key).Bytes()

	if errors.Is(err, redis.Nil) {
		namedLists = append(namedLists, data.NamedList{
			Name:  "",
			Items: []data.ListItem{},
		})
		return namedLists, nil
	} else if err != nil {
		return namedLists, err
	}

	err = json.Unmarshal(jsonData, &namedLists)
	if err != nil {
		return namedLists, err
	}
	return namedLists, nil
}

func (p *Redis) Commit(data []data.NamedList) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	cmd := p.client.Set(ctx, p.key, string(jsonData), 0)

	return cmd.Err()
}
