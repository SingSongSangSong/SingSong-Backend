package config

type RedisConfig struct {
	Addr     string
	Password string // no password set
	DB       int    // use default DB
}

func NewRedisConfig() *RedisConfig {
	return &RedisConfig{Addr: "localhost:6379", Password: "", DB: 0}
}
