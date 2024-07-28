package router

import (
	"SingSong-Server/internal/handler"
	"SingSong-Server/middleware"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/pinecone-io/go-pinecone/pinecone"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

func SetupRouter(db *sql.DB, rdb *redis.Client, idxConnection *pinecone.IndexConnection) *gin.Engine {
	r := gin.Default()

	// CORS 설정 추가
	r.Use(middleware.CORSMiddleware())

	// 추천 엔드포인트 설정
	recommend := r.Group("/api/v1/recommend")
	{
		recommend.POST("/home", handler.HomeRecommendation(db, rdb, idxConnection))
		recommend.POST("/songs", handler.SongRecommendation(db, rdb, idxConnection))
		recommend.POST("/refresh", handler.RefreshRecommendation(rdb, idxConnection))
	}

	// 태그 엔드포인트 설정
	tags := r.Group("/api/v1/tags")
	{
		tags.GET("", handler.ListTags())
	}

	user := r.Group("/api/v1/user")
	{
		user.POST("/login", handler.OAuth(rdb, db))
		user.POST("/reissue", handler.Reissue(rdb))
	}

	// 스웨거 설정
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 404 에러
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "error - invalid api"})
	})

	return r
}
