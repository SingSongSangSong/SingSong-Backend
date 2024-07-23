package model

import (
	"SingSong-Backend/config"
	"context"
	"github.com/redis/go-redis/v9"
	"log"
)

type RedisModel struct {
	redisClient *redis.Client
}

func NewRedisModel(ctx context.Context, config *config.RedisConfig) (*RedisModel, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password, // no password set
		DB:       config.DB,       // use default DB
	})

	// 연결 확인
	if err := rdb.Ping(ctx); err != nil {
		log.Printf("Redis 연결 실패: %v", err)
	}

	return &RedisModel{redisClient: rdb}, nil
}

func (model *RedisModel) Set(ctx context.Context, key string, value interface{}) *redis.StatusCmd {
	return model.redisClient.Set(ctx, key, value, 0)
}

func (model *RedisModel) Get(ctx context.Context, key string) *redis.StringCmd {
	return model.redisClient.Get(ctx, key)
}
