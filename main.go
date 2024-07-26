package main

import (
	"SingSong-Server/conf"
	"SingSong-Server/internal/handler"
	"SingSong-Server/internal/pkg"
	"context"
	"database/sql"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pinecone-io/go-pinecone/pinecone"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"time"

	_ "SingSong-Server/docs"
	_ "github.com/go-sql-driver/mysql"
)

// @title           싱송생송 API
// @version         1.0
// @BasePath  /api/v1
func main() {
	ctx := context.Background()

	var db *sql.DB
	var rdb *redis.Client
	var idxConnection *pinecone.IndexConnection
	conf.SetConfig(ctx, &db, &rdb, &idxConnection)

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
