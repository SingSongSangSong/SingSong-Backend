package main

import (
	"SingSong-Server/conf"
	_ "SingSong-Server/docs"
	"SingSong-Server/router"
	"context"
	"database/sql"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	firebase "firebase.google.com/go/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/pinecone-io/go-pinecone/pinecone"
	"github.com/redis/go-redis/v9"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"
)

// @title           싱송생송 API
// @version         1.0
// @BasePath  /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// 한국 표준시(KST)를 로드하여 전역으로 설정
	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		log.Fatalf("Failed to load location: %v", err)
	}
	time.Local = loc // 서버 전역에서 KST로 처리

	if conf.Env == conf.ProductionMode {
		currentDate := time.Now().Format("2006-01-02")
		gitCommit := os.Getenv("GIT_SHA")
		if gitCommit == "" {
			gitCommit = "unknown" // 기본값 설정, 환경 변수가 없을 경우
		}

		tracer.Start(
			tracer.WithRuntimeMetrics(),
			tracer.WithEnv(conf.Env),
			tracer.WithService(conf.DatadogServiceName),
			tracer.WithServiceVersion(currentDate+":"+gitCommit), //배포날짜:커밋해시로 버전 설정
		)
		defer tracer.Stop()

		err := profiler.Start(
			profiler.WithEnv(conf.Env),
			profiler.WithService(conf.DatadogServiceName),
		)
		if err != nil {
			log.Fatal("Failed to start profiler: ", err)
		}
		defer profiler.Stop()
	}

	ctx := context.Background()

	var db *sql.DB
	var rdb *redis.Client
	var idxConnection *pinecone.IndexConnection
	var milvusClient client.Client
	var s3Client *s3.Client
	var firebaseApp *firebase.App

	conf.SetupConfig(ctx, &db, &rdb, &idxConnection, &milvusClient, &firebaseApp, &s3Client)

	boil.SetDB(db)
	//boil.DebugMode = true

	r := router.SetupRouter(db, rdb, idxConnection, &milvusClient, firebaseApp, s3Client)

	// pprof를 위한 별도의 HTTP 서버 실행
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// 서버 실행
	if err := r.Run(); err != nil {
		log.Fatalf("서버 실행 실패: %v", err)
	}

}
