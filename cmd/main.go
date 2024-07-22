package main

import (
	"SingSong-Backend/config"
	_ "SingSong-Backend/docs"
	"SingSong-Backend/internal/handler"
	"SingSong-Backend/internal/model"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// @title           싱송생송 API
// @version         1.0
// @BasePath  /api/v1

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	ctx := context.Background()

	// MySQL 설정
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatalf("DB_PORT를 숫자로 변환할 수 없습니다: %v", err)
	}

	dbConf := config.NewMysqlConfig(
		os.Getenv("DB_HOST"),
		port,
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"))

	// 모델 초기화
	m, err := model.NewModel(dbConf)
	if err != nil {
		log.Fatalf("Model 생성 실패: %v", err)
	}

	// 핸들러 초기화
	h, err := handler.NewHandler(m)
	if err != nil {
		log.Fatalf("Handler 생성 실패: %v", err)
	}

	// Pinecone 설정
	pineConf := config.NewPineconeConfig(os.Getenv("PINECONE_API_KEY"))

	// Pinecone 클라이언트 초기화
	pc, err := model.NewPineconeClient(ctx, pineConf)
	if err != nil {
		log.Fatalf("Pinecone 생성 실패: %v", err)
	}

	// Pinecone 핸들러 초기화
	ph, err := handler.NewPineconeHandler(pc)
	if err != nil {
		log.Fatalf("NewPineconeHandler 생성 실패")
	}

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

	// 사용자 엔드포인트 설정
	user := r.Group("/user")
	{
		user.GET("", h.ListUser)
		user.POST("", h.RegisterUser)
		user.GET("/:user", h.GetUser)
	}

	// 추천 엔드포인트 설정
	recommend := r.Group("/api/v1/recommend")
	{
		recommend.POST("", ph.RecommendBySongs)
		recommend.POST("/tags", ph.HomeRecommendation)
	}

	// 태그 엔드포인트 설정
	tags := r.Group("/api/v1/tags")
	{
		tags.GET("/ssss", h.ListSsssTags)
	}

	// 스웨거 설정
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 404 에러
	r.NoRoute(func(c *gin.Context) {
		handler.BaseResponse(c, http.StatusNotFound, "error - invalid api", nil)
		return
	})

	// 서버 실행
	if err := r.Run(); err != nil {
		log.Fatalf("서버 실행 실패: %v", err)
	}
}
