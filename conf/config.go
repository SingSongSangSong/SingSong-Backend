package conf

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/pinecone-io/go-pinecone/pinecone"
	"github.com/redis/go-redis/v9"
	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
	"log"
	"os"
	"strconv"
)

type GrpcConfig struct {
	Addr string
}

type AuthConfig struct {
	SECRET_KEY                   string
	KAKAO_REST_API_KEY           string
	KAKAO_ISSUER                 string
	APPLE_CLIENT_ID              string
	APPLE_ISSUER                 string
	JWT_ISSUER                   string
	JWT_ACCESS_VALIDITY_SECONDS  string
	JWT_REFRESH_VALIDITY_SECONDS string
}

type VectorDBConfig struct {
	MILVUS_HOST     string
	MILVUS_PORT     string
	DIMENSION       int
	COLLECTION_NAME string
}

const (
	LocalMode          = "local"
	TestMode           = "test"
	ProductionMode     = "prod"
	DatadogServiceName = "singsong-golang"
)

var (
	AuthConfigInstance     *AuthConfig
	VectorDBConfigInstance *VectorDBConfig
	GrpcConfigInstance     *GrpcConfig
	Env                    string
)

func init() {
	Env = os.Getenv("SERVER_MODE")
	if Env == "" {
		Env = LocalMode // default: local mode
	}

	// 만약 dev면 .env 파일 로드 시도
	if Env == LocalMode {
		log.Println("current environment is local, start to load .env file.")
		err := godotenv.Load(".env")
		if err != nil {
			log.Printf("Error loading .env file during auth configuration.")
		}
	} else if Env == TestMode {
		log.Println("Running in test mode, skip .env file loading.")
	} else {
		log.Println("Running in production mode, skip .env file loading.")
	}

	AuthConfigInstance = &AuthConfig{
		SECRET_KEY:                   os.Getenv("SECRET_KEY"),
		KAKAO_REST_API_KEY:           os.Getenv("KAKAO_REST_API_KEY"),
		KAKAO_ISSUER:                 os.Getenv("KAKAO_ISSUER"),
		APPLE_CLIENT_ID:              os.Getenv("APPLE_CLIENT_ID"),
		APPLE_ISSUER:                 os.Getenv("APPLE_ISSUER"),
		JWT_ISSUER:                   os.Getenv("JWT_ISSUER"),
		JWT_ACCESS_VALIDITY_SECONDS:  os.Getenv("JWT_ACCESS_VALIDITY_SECONDS"),
		JWT_REFRESH_VALIDITY_SECONDS: os.Getenv("JWT_REFRESH_VALIDITY_SECONDS"),
	}

	dimensionStr := os.Getenv("MILVUS_DIMENSION")
	log.Printf("MILVUS_DIMENSION: %s", dimensionStr)
	dimension, err := strconv.Atoi(dimensionStr)
	if err != nil {
		log.Fatalf("Failed to convert DIMENSION to int: %v", err)
	}

	VectorDBConfigInstance = &VectorDBConfig{
		MILVUS_HOST:     os.Getenv("MILVUS_HOST"),
		MILVUS_PORT:     os.Getenv("MILVUS_PORT"),
		DIMENSION:       dimension,
		COLLECTION_NAME: os.Getenv("MILVUS_COLLECTION_NAME"),
	}

	GrpcConfigInstance = &GrpcConfig{
		Addr: func() string {
			if addr := os.Getenv("GRPC_ADDR"); addr != "" {
				return addr
			}
			return "python-gRPC" // 기본값
		}(),
	}
}

func SetupConfig(ctx context.Context, db **sql.DB, rdb **redis.Client, idxConnection **pinecone.IndexConnection, milvusClient *client.Client) {
	var err error

	// MySQL 연결 설정
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	sqltrace.Register("mysql", &mysql.MySQLDriver{}, sqltrace.WithServiceName("singsong-mysql"))
	*db, err = sqltrace.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Mysql 연결 실패: %v", err)
	}
	if err := (*db).Ping(); err != nil {
		log.Fatalf("Mysql ping 실패: %v", err)
	}

	// 레디스
	*rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"), // no password set
	})

	_, err = (*rdb).Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Redis 연결 실패: %v", err)
	}

	// Milvus 연결
	*milvusClient, err = client.NewClient(ctx, client.Config{Address: os.Getenv("MILVUS_HOST") + ":" + os.Getenv("MILVUS_PORT")})
	if err != nil {
		log.Fatalf("Milvus 연결 실패: %v", err)
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
