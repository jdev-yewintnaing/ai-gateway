package ratelimit

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Limiter struct {
	client *redis.Client
	limit  int
}

func NewLimiter(redisURL string, limit int) (*Limiter, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opts)
	return &Limiter{client: client, limit: limit}, nil
}

func (l *Limiter) Allow(ctx context.Context, caller string, tokens int) (bool, error) {
	if l.client == nil {
		return true, nil
	}

	now := time.Now()
	key := fmt.Sprintf("rl:tokens:%s:%s", caller, now.Format("200601021504"))

	// Atomic check and increment
	res, err := IncrementAndCheckLua.Run(ctx, l.client, []string{key}, tokens, l.limit, 120).Result()
	if err != nil {
		return false, err
	}

	allowed := res.(int64) == 1
	return allowed, nil
}
