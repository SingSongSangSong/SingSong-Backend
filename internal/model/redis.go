package model

import (
	"SingSong-Backend/config"
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
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
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Printf("Redis 연결 실패: %v", err)
		return nil, err
	}
	log.Println("Redis 연결 성공:", pong)

	return &RedisModel{redisClient: rdb}, nil
}

func (model *RedisModel) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return model.redisClient.Set(ctx, key, value, expiration)
}

func (model *RedisModel) Get(ctx context.Context, key string) *redis.StringCmd {
	return model.redisClient.Get(ctx, key)
}

func (model *RedisModel) SavePublicKey(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return model.redisClient.Set(ctx, key, value, expiration)
}
