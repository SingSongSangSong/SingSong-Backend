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
	_ "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
	"net/http"
)

func SetupRouter(db *sql.DB, rdb *redis.Client, idxConnection *pinecone.IndexConnection) *gin.Engine {
	// Initialize Datadog tracer
	//tracer.Start()
	//defer tracer.Stop()

	r := gin.Default()

	// Wrap the router with Datadog middleware
	//r.Use(ddgin.Middleware("singsong-service"))

	// CORS 설정 추가
	r.Use(middleware.PlatformMiddleware())

	// 버전 확인
	version := r.Group("/api/v1/version")
	{
		version.GET("/", handler.AllVersion(db))
		version.POST("/check", middleware.PlatformMiddleware(), handler.VersionCheck(db))
		version.POST("/update", handler.LatestVersionUpdate(db))
	}

	r.GET("/", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "Welcome to SingSong-Server"}) })

	// 추천 엔드포인트 설정
	recommend := r.Group("/api/v1/recommend")
	//recommend.Use(middleware.AuthMiddleware()) // 추천 엔드포인트 전체에서 인증을 쓴다면 이렇게도 가능
	{
		recommend.POST("/home", handler.HomeRecommendation(db, rdb, idxConnection))
		recommend.GET("/home/songs", handler.HomeSongRecommendation(db))
		recommend.POST("/songs", handler.SongRecommendation(db, rdb, idxConnection))
		recommend.POST("/refresh", middleware.AuthMiddleware(db), handler.RefreshRecommendation(db, rdb, idxConnection)) //일단 새로고침에만 적용
	}

	// 태그 엔드포인트 설정
	tags := r.Group("/api/v1/tags")
	{
		tags.GET("", handler.ListTags())
	}

	member := r.Group("/api/v1/member")
	{
		member.POST("/login", handler.Login(rdb, db))
		member.POST("/reissue", handler.Reissue(rdb))
		member.GET("", middleware.AuthMiddleware(db), handler.GetMemberInfo(db))
		member.POST("/withdraw", middleware.AuthMiddleware(db), handler.Withdraw(db, rdb))
	}

	// 태그 엔드포인트 설정
	keep := r.Group("/api/v1/keep")
	{
		keep.GET("", middleware.AuthMiddleware(db), handler.GetSongsFromPlaylist(db))
		keep.POST("", middleware.AuthMiddleware(db), handler.AddSongsToPlaylist(db))
		keep.DELETE("", middleware.AuthMiddleware(db), handler.DeleteSongsFromPlaylist(db))
	}

	// 스웨거 설정
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 404 에러
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "error - invalid api"})
	})

	return r
}
