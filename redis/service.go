package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisConf struct {
	Address    string
	ExpireTime time.Duration
	Size       int
}

type Client struct {
	RedisClient *redis.Client
}

func NewClient(c RedisConf) Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Address,
		Password: "",
		DB:       0,
	})

	log.Print("Setting up redis connection")

	if err := rdb.Ping(context.TODO()).Err(); err != nil {
		log.Fatal(err)
	}

	cli := Client{
		RedisClient: rdb,
	}

	return cli
}

func (c Client) Set(ctx context.Context, key, val string) error {
	// err := c.RedisClient.Set(ctx, key, val, 0).Err()

	pipe := c.RedisClient.TxPipeline()
	pipe.Set(ctx, key, val, 0)
	_, err := pipe.Exec(ctx)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (c Client) Get(ctx context.Context, key string) (string, error) {
	val, err := c.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return val, err
	}

	return val, nil
}

func (c Client) GetAll(ctx context.Context) ([]string, error) {
	var vals []string
	iter := c.RedisClient.Scan(ctx, 0, "key*", 10).Iterator()

	for iter.Next(ctx) {
		fmt.Println(iter.Val())
		vals = append(vals, iter.Val())
	}

	if err := iter.Err(); err != nil {
		fmt.Println(err)
		return vals, err
	}

	return vals, nil
}
