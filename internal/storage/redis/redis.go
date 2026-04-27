package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/north-fy/golang-restapi-todolist/internal/config"
	"github.com/redis/go-redis/v9"
)

const op = "storage/redis/"

type Storage struct {
	client *redis.Client
	ttl    time.Duration
}

func NewStorage(cfg config.RedisConfig) *Storage {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		Username: cfg.User,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}

	panic("123")

	return &Storage{
		client: client,
		ttl:    cfg.TTL,
	}
}

func (s *Storage) Set(ctx context.Context, key string, value any) error {
	return s.client.Set(ctx, key, value, s.ttl).Err()
}

func (s *Storage) Get(ctx context.Context, key string, result interface{}) error {
	res := s.client.Get(ctx, key)
	if err := res.Err(); err != nil {
		return err
	}

	return res.Scan(result)
}

func (s *Storage) Close() error {
	return s.client.Close()
}
