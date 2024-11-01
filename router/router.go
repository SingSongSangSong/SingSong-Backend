package router

import (
	"SingSong-Server/conf"
	"SingSong-Server/internal/handler"
	"SingSong-Server/middleware"
	"database/sql"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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

func SetupRouter(db *sql.DB, rdb *redis.Client, idxConnection *pinecone.IndexConnection, milvusClient *client.Client, s3Client *s3.Client) *gin.Engine {

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
		version.POST("/check", handler.VersionCheck(db))
		version.POST("/update", handler.VersionUpdate(db))
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

	recommendV2 := r.Group("/api/v2/recommend")
	{
		recommendV2.GET("recommendation/ai", middleware.AuthMiddleware(db), handler.GetRecommendationV2(db, rdb, milvusClient))
		recommendV2.GET("recommendation/:pageId", middleware.AuthMiddleware(db), handler.GetRecommendationV2(db, rdb, milvusClient))
		recommendV2.POST("/refresh", middleware.AuthMiddleware(db), handler.RefreshRecommendationV2(db))
		recommendV2.POST("/recommendation/functionCallingWithTypes", middleware.AuthMiddleware(db), handler.FunctionCallingWithTypesRecommedation(db))
		recommendV2.GET("/recommendation/searchLog", middleware.AuthMiddleware(db), handler.GetSearchResultsForLLMV2(db))
	}

	// 태그 엔드포인트 설정
	tags := r.Group("/api/v1/tags")
	{
		tags.GET("", handler.ListTags())
	}

	tagsV2 := r.Group("/api/v2/tags")
	{
		tagsV2.GET("", handler.ListTagsV2())
	}

	tagsV3 := r.Group("/api/v3/tags")
	{
		tagsV3.GET("", handler.ListTagsV3())
	}

	tagsV4 := r.Group("/api/v4/tags")
	{
		tagsV4.GET("", handler.ListTagsV4())
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
		keep.GET("", middleware.AuthMiddleware(db), handler.GetSongsFromKeep(db))
		keep.POST("", middleware.AuthMiddleware(db), handler.AddSongsToKeep(db))
		keep.DELETE("", middleware.AuthMiddleware(db), handler.DeleteSongsFromKeep(db))
		keep.GET("/story", middleware.AuthMiddleware(db), handler.GetKeepForStory(db))
		keep.POST("/:keepListId/like", middleware.AuthMiddleware(db), handler.KeepListLike(db))
		keep.GET("/:keepListId", middleware.AuthMiddleware(db), handler.GetSongsFromKeepInStory(db))
		keep.POST("/:keepListId/subscribe", middleware.AuthMiddleware(db), handler.SubscribeKeep(db))
		keep.GET("/subscribe", middleware.AuthMiddleware(db), handler.GetSubscribedKeeps(db))
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
		songs.GET("/:songId/comments/hot", middleware.AuthMiddleware(db), handler.GetHotCommentOfSong(db))
	}

	songsV2 := r.Group("/api/v2/songs")
	{
		songsV2.GET("/:songId/related", middleware.AuthMiddleware(db), handler.RelatedSongV2(db, milvusClient))
		songsV2.GET("/:songId/comments", middleware.AuthMiddleware(db), handler.GetCommentsOnSongV2(db))
	}

	songV3 := r.Group("/api/v3/songs")
	{
		songV3.GET("/:songId/comments", middleware.AuthMiddleware(db), handler.GetCommentsOnSongV3(db))
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
		comment.DELETE("/:commentId", middleware.AuthMiddleware(db), handler.DeleteComment(db))
		comment.GET("/my", middleware.AuthMiddleware(db), handler.GetMySongComment(db))
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
		search.GET("/posts", middleware.AuthMiddleware(db), handler.SearchPosts(db))
	}

	searchV2 := r.Group("/api/v2/search")
	{
		searchV2.GET("/:searchKeyword", middleware.AuthMiddleware(db), handler.SearchSongsV2(db))
		searchV2.GET("/artist-name", middleware.AuthMiddleware(db), handler.SearchSongsByAristV2(db))
		searchV2.GET("/song-name", middleware.AuthMiddleware(db), handler.SearchSongsBySongNameV2(db))
		searchV2.GET("/song-number", middleware.AuthMiddleware(db), handler.SearchSongsBySongNumberV2(db))
	}

	post := r.Group("/api/v1/posts")
	{
		post.POST("", middleware.AuthMiddleware(db), handler.CreatePost(db))
		post.GET("", middleware.AuthMiddleware(db), handler.ListPosts(db))
		post.GET("/:postId", middleware.AuthMiddleware(db), handler.GetPost(db))
		post.DELETE("/:postId", middleware.AuthMiddleware(db), handler.DeletePost(db))
		post.POST("/:postId/reports", middleware.AuthMiddleware(db), handler.ReportPost(db))
		post.POST("/:postId/likes", middleware.AuthMiddleware(db), handler.LikePost(db))
		post.GET("/:postId/comments", middleware.AuthMiddleware(db), handler.GetCommentOnPost(db))
	}

	postV2 := r.Group("/api/v2/posts")
	{
		postV2.GET("/:postId/comments", middleware.AuthMiddleware(db), handler.GetCommentOnPostV2(db))
	}

	postComment := r.Group("/api/v1/posts/comments")
	{
		postComment.POST("", middleware.AuthMiddleware(db), handler.CommentOnPost(db))
		postComment.GET("/:postCommentId/recomments", middleware.AuthMiddleware(db), handler.GetReCommentOnPost(db))
		postComment.POST("/report", middleware.AuthMiddleware(db), handler.ReportPostComment(db))
		postComment.POST("/:postCommentId/like", middleware.AuthMiddleware(db), handler.LikePostComment(db))
		postComment.DELETE("/:postCommentId", middleware.AuthMiddleware(db), handler.DeletePostComment(db))
	}

	recent := r.Group("/api/v1/recent")
	{
		recent.GET("/search", handler.GetLatestSearchApi(db))
		recent.GET("/keep", handler.GetRecentKeepSongs(db))
		recent.GET("/comment", handler.GetRecentCommentsongs(db))
	}

	record := r.Group("/api/v1/record")
	{
		record.POST("/song", middleware.AuthMiddleware(db), handler.RecordSong(db, s3Client))
	}

	// 스웨거 설정
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 404 에러
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "error - invalid api"})
	})

	return r
}
