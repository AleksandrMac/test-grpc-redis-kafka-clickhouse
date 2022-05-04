package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/AleksandrMac/test_hezzl/pkg/user/model"
	redisLib "github.com/go-redis/redis/v8"
)

type RedisClient struct {
	client *redisLib.Client
	ttl    time.Duration
}

func New(client *redisLib.Client, ttl time.Duration) *RedisClient {
	return &RedisClient{
		client: client,
		ttl:    ttl,
	}
}

func (x *RedisClient) GetList(ctx context.Context, limit, offset uint64) ([]model.User, error) {
	data, err := x.client.Get(context.Background(), fmt.Sprintf("userList_l%d.o%d", limit, offset)).Bytes()
	if err == redisLib.Nil {
		// we got empty result, it's not an error
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var users []model.User
	err = json.Unmarshal(data, &users)
	return users, err
}

func (x *RedisClient) SetList(ctx context.Context, limit, offset uint64, users []model.User) error {
	data, err := json.Marshal(users)
	if err != nil {
		return err
	}

	return x.client.Set(
		ctx,
		fmt.Sprintf("userList_l%d.o%d", limit, offset),
		data,
		x.ttl,
	).Err()
}
