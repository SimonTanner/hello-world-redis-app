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

	cli := Client{
		RedisClient: rdb,
		expireTime:  c.ExpireTime,
	}

	return cli
}

func (c *Client) Set(ctx context.Context, key, val string) error {
	// err := c.RedisClient.Set(ctx, key, val, 0).Err()

	pipe := c.RedisClient.TxPipeline()
	pipe.Set(ctx, key, val, c.expireTime)
	_, err := pipe.Exec(ctx)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	val, err := c.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return val, err
	}

	return val, nil
}

func (c *Client) GetAll(ctx context.Context) ([]string, error) {
	var vals []string
	var cursor uint64
	for {
		var (
			keys []string
			err  error
		)
		keys, cursor, err = c.RedisClient.Scan(ctx, cursor, "*", 0).Result()
		if err != nil {
			return vals, err
		}

		for idx, key := range keys {
			fmt.Println(idx, key)
			val, err := c.RedisClient.Get(ctx, key).Result()
			if err != nil {
				return vals, err
			}
			fmt.Println("CURSOR:", cursor, "KEY:", key, "VALUES:", val)

			vals = append(vals, val)
			if idx == len(keys)-1 {
				return vals, nil
			}
		}

		if cursor == 0 {
			break
		}
	}

	return vals, nil
}

type Object struct {
	Str string
	Num int
}

func CheckCache() {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server1": ":6379",
		},
	})

	mycache := cache.New(&cache.Options{
		Redis:      ring,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	ctx := context.TODO()
	key := "mykey"
	obj := &Object{
		Str: "mystring",
		Num: 42,
	}

	if err := mycache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: obj,
		TTL:   time.Hour,
	}); err != nil {
		panic(err)
	}

	var wanted Object
	if err := mycache.Get(ctx, key, &wanted); err == nil {
		fmt.Println("wanted:", wanted)
	}
}
