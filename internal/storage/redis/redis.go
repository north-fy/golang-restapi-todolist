package redis

import (
	"context"
	"encoding/json"
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
	opt, err := redis.ParseURL(fmt.Sprintf("redis://%s:%s@%s:%d/%d", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB))
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(opt)

	return &Storage{
		client: client,
		ttl:    cfg.TTL,
	}
}

func (s *Storage) Set(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return s.client.Set(ctx, key, data, s.ttl).Err()
}

func (s *Storage) Get(ctx context.Context, key string, obj interface{}) error {
	res := s.client.Get(ctx, key)
	if err := res.Scan(obj); err != nil {
		return err
	}

	return nil
}

func (s *Storage) Close() error {
	return s.client.Close()
}
