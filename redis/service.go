package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

type RedisConf struct {
	Address    string
	ExpireTime time.Duration
	Size       int
}

type Client struct {
	RedisClient *redis.Client
	RedisCache  *cache.Cache
	expireTime  time.Duration
}

func NewClient(c RedisConf) Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Address,
		Password: "",
		DB:       0,
	})

	log.Print("Checking connection to redis-server")

	if err := rdb.Ping(context.TODO()).Err(); err != nil {
		log.Fatal(err)
	}

	rCache := cache.New(&cache.Options{
		Redis:      rdb,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	cli := Client{
		RedisClient: rdb,
		RedisCache:  rCache,
		expireTime:  c.ExpireTime,
	}

	return cli
}

type Message struct {
	Key string
	Str string
}

func (c *Client) Set(ctx context.Context, key string, msg Message) error {
	if err := c.RedisCache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: msg,
		TTL:   c.expireTime,
	}); err != nil {
		return err
	}

	return nil
}

func (c *Client) Get(ctx context.Context, key string) (Message, error) {
	var msg Message
	err := c.RedisCache.Get(ctx, key, &msg)
	if err != nil {
		return msg, err
	}
	msg.Key = key

	return msg, nil
}

func (c *Client) GetAll(ctx context.Context) ([]Message, error) {
	var msgs []Message
	iter := c.RedisClient.Scan(ctx, 0, "*", 0).Iterator()
	for iter.Next(ctx) {
		fmt.Println("Key:", iter.Val())
		key := iter.Val()
		msg, err := c.Get(ctx, key)
		if err != nil {
			return msgs, err
		}
		msgs = append(msgs, msg)
	}

	return msgs, iter.Err()
}
