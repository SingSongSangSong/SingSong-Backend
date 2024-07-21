package main

import (
	"SingSong-Backend/config"
	_ "SingSong-Backend/docs"
	"SingSong-Backend/internal/handler"
	"SingSong-Backend/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/net/context"
	"log"
	"os"
	"strconv"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	err := godotenv.Load(".env")
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
		recommend.POST("", ph.RegisterRecommendation)
		recommend.POST("/tags", ph.HomeRecommendation)
	}

	//스웨거 설정
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 서버 실행
	if err := r.Run(); err != nil {
		log.Fatalf("서버 실행 실패: %v", err)
	}
}
