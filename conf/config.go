package conf

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/pinecone-io/go-pinecone/pinecone"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
)

type AuthConfig struct {
	SECRET_KEY                   string
	KAKAO_REST_API_KEY           string
	KAKAO_ISSUER                 string
	JWT_ISSUER                   string
	JWT_ACCESS_VALIDITY_SECONDS  string
	JWT_REFRESH_VALIDITY_SECONDS string
}

var (
	AuthConfigInstance *AuthConfig
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Error loading .env file during auth configuration. ") //개발환경용
	}
	AuthConfigInstance = &AuthConfig{
		SECRET_KEY:                   os.Getenv("SECRET_KEY"),
		KAKAO_REST_API_KEY:           os.Getenv("KAKAO_REST_API_KEY"),
		KAKAO_ISSUER:                 os.Getenv("KAKAO_ISSUER"),
		JWT_ISSUER:                   os.Getenv("JWT_ISSUER"),
		JWT_ACCESS_VALIDITY_SECONDS:  os.Getenv("JWT_ACCESS_VALIDITY_SECONDS"),
		JWT_REFRESH_VALIDITY_SECONDS: os.Getenv("JWT_REFRESH_VALIDITY_SECONDS"),
	}
}

func SetupConfig(ctx context.Context, db **sql.DB, rdb **redis.Client, idxConnection **pinecone.IndexConnection) {
	var err error
	// MySQL 설정
	err = godotenv.Load(".env")
	if err != nil {
		log.Printf("Error loading .env file") //개발환경용
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	*db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Mysql 연결 실패: %v", err)
	}

	if err := (*db).Ping(); err != nil {
		log.Fatalf("Mysql ping 실패: %v", err)
	}

	// 레디스
	*rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"), // no password set
	})

	_, err = (*rdb).Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Redis 연결 실패: %v", err)
	}

	// Pinecone 연결
	pineconeApiKey := os.Getenv("PINECONE_API_KEY")
	if pineconeApiKey == "" {
		log.Fatalf("Pinecone api key 없음")
	}

	client, err := pinecone.NewClient(pinecone.NewClientParams{ApiKey: pineconeApiKey})
	if err != nil {
		log.Fatalf("Pinecone 실패: %v", err)
	}

	idx, err := client.DescribeIndex(ctx, os.Getenv("PINECONE_INDEX"))
	if err != nil {
		log.Fatalf("Failed to describe index \"%s\". Error:%s", idx.Name, err)
	}

	*idxConnection, err = client.Index(idx.Host)
	if err != nil {
		log.Fatalf("Failed to create IndexConnection for Host: %v. Error: %v", idx.Host, err)
	}
}
