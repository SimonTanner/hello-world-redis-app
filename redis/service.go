package redis

import (
	"context"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

type RedisConf struct {
	Address    map[string]string
	ExpireTime time.Duration
	Size       int
}

type Client struct {
	RedisCache *cache.Cache
}

func NewClient(c RedisConf) Client {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: c.Address,
	})

	rcache := cache.New(&cache.Options{
		Redis:      ring,
		LocalCache: cache.NewTinyLFU(c.Size, c.ExpireTime),
	})

	cli := Client{
		RedisCache: rcache,
	}

	return cli
}

func (c Client) Set(ctx context.Context, key, val string) error {
	err := c.RedisCache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: val,
	})

	if err != nil {
		return err
	}

	return nil
}

func (c Client) Get(ctx context.Context, key string) (string, error) {
	var val string
	err := c.RedisCache.Get(ctx, key, &val)
	if err != nil {
		return val, err
	}

	return val, nil
}
