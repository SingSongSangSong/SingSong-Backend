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
)

// @title           싱송생송 API
// @version         1.0
// @BasePath  /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	if conf.Env == conf.ProductionMode {
		tracer.Start(
			tracer.WithRuntimeMetrics(),
			tracer.WithEnv(conf.Env),
			tracer.WithService("singsong"),
			tracer.WithServiceVersion("2024.08.29"), // todo: 버전 수정 자동으로
			tracer.WithDebugMode(true),
			tracer.WithAnalytics(true),
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
