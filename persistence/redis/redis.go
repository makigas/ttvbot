package redis

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"gopkg.makigas.es/ttvbot/domain/config"
)

var Module = fx.Module("Redis", fx.Provide(NewRedisClient))

type RedisClient struct {
	fx.Out
	Client *redis.Client
}

func NewRedisClient(cfg *config.Config) RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Database,
	})
	return RedisClient{Client: client}
}
