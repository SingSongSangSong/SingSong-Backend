package main

import (
	"SingSong-Server/internal/handler"
	"SingSong-Server/internal/pkg"
	"context"
	"database/sql"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pinecone-io/go-pinecone/pinecone"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"os"
	"time"

	_ "SingSong-Server/docs"
	_ "github.com/go-sql-driver/mysql"
)

func setConfig(ctx context.Context, db **sql.DB, rdb **redis.Client, idxConnection **pinecone.IndexConnection) {
	var err error
	// MySQL 설정
	err = godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file") //개발환경용
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

// @title           싱송생송 API
// @version         1.0
// @BasePath  /api/v1
func main() {
	ctx := context.Background()

	var db *sql.DB
	var rdb *redis.Client
	var idxConnection *pinecone.IndexConnection
	setConfig(ctx, &db, &rdb, &idxConnection)

	// Gin 라우터 설정
	r := gin.Default()

	// CORS 설정 추가
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 추천 엔드포인트 설정
	recommend := r.Group("/api/v1/recommend")
	{
		recommend.POST("/home", handler.HomeRecommendation(db, rdb, idxConnection))
		recommend.POST("/songs", handler.SongRecommendation(db, rdb, idxConnection))
	}

	// 태그 엔드포인트 설정
	tags := r.Group("/api/v1/tags")
	{
		tags.GET("", handler.ListTags())
	}

	// 스웨거 설정
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 404 에러
	r.NoRoute(func(c *gin.Context) {
		pkg.BaseResponse(c, http.StatusNotFound, "error - invalid api", nil)
		return
	})

	// 서버 실행
	if err := r.Run(); err != nil {
		log.Fatalf("서버 실행 실패: %v", err)
	}
}
