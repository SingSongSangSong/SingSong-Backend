package router

import (
	"SingSong-Server/conf"
	"SingSong-Server/internal/handler"
	"SingSong-Server/middleware"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/pinecone-io/go-pinecone/pinecone"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
	"net/http"
)

func SetupRouter(db *sql.DB, rdb *redis.Client, idxConnection *pinecone.IndexConnection) *gin.Engine {

	//gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Datadog tracer
	if conf.Env == conf.ProductionMode {
		r.Use(gintrace.Middleware("singsong"))
	}

	// CORS 설정 추가
	r.Use(middleware.CORSMiddleware())

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
		recommend.POST("/refresh", middleware.AuthMiddleware(db), handler.RefreshRecommendation(db, rdb, idxConnection)) //일단 새로고침에만 적용
		recommend.GET("/recommendation", middleware.AuthMiddleware(db), handler.GetRecommendation(db, rdb, idxConnection))
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
		member.POST("/logout", middleware.AuthMiddleware(db), handler.Logout(rdb))
	}

	// 태그 엔드포인트 설정
	keep := r.Group("/api/v1/keep")
	{
		keep.GET("", middleware.AuthMiddleware(db), handler.GetSongsFromPlaylist(db))
		keep.POST("", middleware.AuthMiddleware(db), handler.AddSongsToPlaylist(db))
		keep.DELETE("", middleware.AuthMiddleware(db), handler.DeleteSongsFromPlaylist(db))
	}

	// 노래 상세
	songs := r.Group("/api/v1/songs")
	{
		songs.GET("/:songId", middleware.AuthMiddleware(db), handler.GetSongInfo(db))
		songs.GET("/:songId/reviews", middleware.AuthMiddleware(db), handler.GetSongReview(db))
		songs.PUT("/:songId/reviews", middleware.AuthMiddleware(db), handler.PutSongReview(db))
		songs.DELETE("/:songId/reviews", middleware.AuthMiddleware(db), handler.DeleteSongReview(db))
		songs.GET("/:songId/related", middleware.AuthMiddleware(db), handler.RelatedSong(db, idxConnection))
	}

	// 노래 리뷰 선택지 추가/조회
	songReviewOptions := r.Group("/api/v1/song-review-options")
	{
		songReviewOptions.GET("", handler.ListSongReviewOptions(db))
		songReviewOptions.POST("", handler.AddSongReviewOption(db))
	}

	comment := r.Group("/api/v1/comment")
	{
		comment.POST("", middleware.AuthMiddleware(db), handler.CommentOnSong(db))
		comment.GET("/:songId", middleware.AuthMiddleware(db), handler.GetCommentOnSong(db))
		comment.POST("/report", middleware.AuthMiddleware(db), handler.ReportComment(db))
		comment.GET("/recomment/:commentId", middleware.AuthMiddleware(db), handler.GetReCommentOnSong(db))
		comment.POST("/:commentId/like", middleware.AuthMiddleware(db), handler.LikeComment(db))
	}

	blacklist := r.Group("/api/v1/blacklist")
	{
		blacklist.POST("", middleware.AuthMiddleware(db), handler.AddBlacklist(db))
		blacklist.DELETE("", middleware.AuthMiddleware(db), handler.DeleteBlacklist(db))
		blacklist.GET("", middleware.AuthMiddleware(db), handler.GetBlacklist(db))
	}

	chart := r.Group("/api/v1/chart")
	{
		chart.GET("", middleware.AuthMiddleware(db), handler.GetChart(rdb))
	}

	search := r.Group("/api/v1/search")
	{
		search.GET("/:searchKeyword", handler.SearchSongs(db))
		search.GET("/artist", handler.SearchSongsByArist(db))
		search.GET("/song-name", handler.SearchSongsBySongName(db))
		search.GET("/song-number", handler.SearchSongsBySongNumber(db))
	}

	// 스웨거 설정
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 404 에러
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "error - invalid api"})
	})

	return r
}
