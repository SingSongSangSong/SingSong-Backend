package router

import (
	"SingSong-Server/conf"
	"SingSong-Server/internal/handler"
	"SingSong-Server/middleware"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/pinecone-io/go-pinecone/pinecone"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
	"net/http"
)

func SetupRouter(db *sql.DB, rdb *redis.Client, idxConnection *pinecone.IndexConnection, milvusClient *client.Client) *gin.Engine {

	//gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Datadog tracer
	if conf.Env == conf.ProductionMode {
		r.Use(gintrace.Middleware(conf.DatadogServiceName))
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
		recommend.GET("/recommendation/ai", middleware.AuthMiddleware(db), handler.GetRecommendation(db, rdb))
		recommend.GET("/recommendation/:pageId", middleware.AuthMiddleware(db), handler.GetRecommendation(db, rdb))
		recommend.POST("/recommendation/llm", middleware.AuthMiddleware(db), handler.LlmHandler(db))
		recommend.POST("/recommendation/langchainAgent", middleware.AuthMiddleware(db), handler.LangchainAgentRecommedation(db))
		recommend.POST("/recommendation/functionCalling", middleware.AuthMiddleware(db), handler.FunctionCallingRecommedation(db))
		recommend.GET("/recommendation/searchLog", middleware.AuthMiddleware(db), handler.GetSearchResultsForLLM(db))
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
		member.PATCH("/nickname", middleware.AuthMiddleware(db), handler.UpdateNickname(db))
	}

	memberV2 := r.Group("/api/v2/member")
	{
		memberV2.POST("/login", handler.LoginV2(rdb, db))
		memberV2.POST("/login/extra", middleware.AuthMiddleware(db), handler.LoginV2ExtraInfoRequired(db))
	}

	// 태그 엔드포인트 설정
	keep := r.Group("/api/v1/keep")
	{
		keep.GET("", middleware.AuthMiddleware(db), handler.GetSongsFromPlaylist(db))
		keep.POST("", middleware.AuthMiddleware(db), handler.AddSongsToPlaylist(db))
		keep.DELETE("", middleware.AuthMiddleware(db), handler.DeleteSongsFromPlaylist(db))
	}

	keepV2 := r.Group("/api/v2/keep")
	{
		keepV2.GET("", middleware.AuthMiddleware(db), handler.GetSongsFromPlaylistV2(db))
	}

	// 노래 상세
	songs := r.Group("/api/v1/songs")
	{
		songs.GET("/:songId", middleware.AuthMiddleware(db), handler.GetSongInfo(db))
		songs.GET("/:songId/reviews", middleware.AuthMiddleware(db), handler.GetSongReview(db))
		songs.PUT("/:songId/reviews", middleware.AuthMiddleware(db), handler.PutSongReview(db))
		songs.DELETE("/:songId/reviews", middleware.AuthMiddleware(db), handler.DeleteSongReview(db))
		songs.GET("/:songId/related", middleware.AuthMiddleware(db), handler.RelatedSong(db, idxConnection))
		songs.GET("/:songId/link", handler.GetLinkBySongInfoId(db))
		songs.GET("/new", middleware.AuthMiddleware(db), handler.ListNewSongs(db))
	}

	songsV2 := r.Group("/api/v2/songs")
	{
		songsV2.GET("/:songId/related", middleware.AuthMiddleware(db), handler.RelatedSongV2(db, milvusClient))
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
		comment.GET("/latest", middleware.AuthMiddleware(db), handler.GetLatestComments(db))
	}

	blacklist := r.Group("/api/v1/blacklist")
	{
		blacklist.POST("", middleware.AuthMiddleware(db), handler.AddBlacklist(db))
		blacklist.DELETE("", middleware.AuthMiddleware(db), handler.DeleteBlacklist(db))
		blacklist.GET("", middleware.AuthMiddleware(db), handler.GetBlacklist(db))
	}

	chartV1 := r.Group("/api/v1/chart")
	{
		chartV1.GET("", middleware.AuthMiddleware(db), handler.GetChart(rdb))
	}

	chartV2 := r.Group("/api/v2/chart")
	{
		chartV2.GET("", middleware.AuthMiddleware(db), handler.GetChartV2(rdb))
	}

	search := r.Group("/api/v1/search")
	{
		search.GET("/:searchKeyword", handler.SearchSongs(db))
		search.GET("/artist-name", handler.SearchSongsByArist(db))
		search.GET("/song-name", handler.SearchSongsBySongName(db))
		search.GET("/song-number", handler.SearchSongsBySongNumber(db))
		search.GET("/posts", handler.SearchPosts(db))
	}

	post := r.Group("/api/v1/posts")
	{
		post.POST("", middleware.AuthMiddleware(db), handler.CreatePost(db))
		post.GET("", handler.ListPosts(db))
		post.GET("/:postId", middleware.AuthMiddleware(db), handler.GetPost(db))
		post.DELETE("/:postId", middleware.AuthMiddleware(db), handler.DeletePost(db))
		post.POST("/:postId/reports", middleware.AuthMiddleware(db), handler.ReportPost(db))
		post.POST("/:postId/likes", middleware.AuthMiddleware(db), handler.LikePost(db))
		post.GET("/:postId/comments", middleware.AuthMiddleware(db), handler.GetCommentOnPost(db))
	}

	postComment := r.Group("/api/v1/posts/comments")
	{
		postComment.POST("", middleware.AuthMiddleware(db), handler.CommentOnPost(db))
		postComment.GET("/:postCommentId/recomments", middleware.AuthMiddleware(db), handler.GetReCommentOnPost(db))
		postComment.POST("/report", middleware.AuthMiddleware(db), handler.ReportPostComment(db))
		postComment.POST("/:postCommentId/like", middleware.AuthMiddleware(db), handler.LikePostComment(db))
	}

	// 스웨거 설정
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 404 에러
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "error - invalid api"})
	})

	return r
}
