package cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewCache(redisURL string, ttl time.Duration) (*Cache, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opts)
	// Test connection
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis for caching: %w", err)
	}
	return &Cache{client: client, ttl: ttl}, nil
}

func (c *Cache) Get(ctx context.Context, key string, target interface{}) (bool, error) {
	if c.client == nil {
		return false, nil
	}

	val, err := c.client.Get(ctx, "cache:"+key).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}

	if err := json.Unmarshal([]byte(val), target); err != nil {
		return false, err
	}

	return true, nil
}

func (c *Cache) Set(ctx context.Context, key string, value interface{}) error {
	if c.client == nil {
		return nil
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, "cache:"+key, string(data), c.ttl).Err()
}

func GenerateKey(model string, messages interface{}) (string, error) {
	data, err := json.Marshal(messages)
	if err != nil {
		return "", err
	}

	h := sha256.New()
	h.Write([]byte(model))
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil)), nil
}
