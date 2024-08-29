package main

import (
	"SingSong-Server/conf"
	_ "SingSong-Server/docs"
	"SingSong-Server/router"
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pinecone-io/go-pinecone/pinecone"
	"github.com/redis/go-redis/v9"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
	"log"
	"os"
	"time"
)

// @title           싱송생송 API
// @version         1.0
// @BasePath  /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	if conf.Env == conf.ProductionMode {
		currentDate := time.Now().Format("2006-01-02")
		gitCommit := os.Getenv("GIT_SHA")
		if gitCommit == "" {
			gitCommit = "unknown" // 기본값 설정, 환경 변수가 없을 경우
		}

		tracer.Start(
			tracer.WithRuntimeMetrics(),
			tracer.WithEnv(conf.Env),
			tracer.WithService("singsong"),
			tracer.WithServiceVersion(currentDate+":"+gitCommit), //배포날짜:커밋해시로 버전 설정
		)
		defer tracer.Stop()

		err := profiler.Start(
			profiler.WithEnv(conf.Env),
			profiler.WithService("singsong"),
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
	conf.SetupConfig(ctx, &db, &rdb, &idxConnection)
	// SQLBoiler의 디버그 모드 활성화
	//boil.DebugMode = true

	r := router.SetupRouter(db, rdb, idxConnection)

	// 서버 실행
	if err := r.Run(); err != nil {
		log.Fatalf("서버 실행 실패: %v", err)
	}

}
