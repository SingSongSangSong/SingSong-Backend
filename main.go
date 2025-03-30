package main

import (
	"SingSong-Server/conf"
	_ "SingSong-Server/docs"
	"SingSong-Server/router"
	"context"
	"database/sql"
	"errors"
	firebase "firebase.google.com/go/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	_ "github.com/go-sql-driver/mysql"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/pinecone-io/go-pinecone/pinecone"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
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
	logrus.SetFormatter(&logrus.JSONFormatter{})
	ctx := context.Background()

	if conf.Env == conf.ProductionMode || conf.Env == conf.TestMode {
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
			logrus.WithContext(ctx).Fatal("Failed to start profiler: ", err)
		}
		defer profiler.Stop()
	}

	var db *sql.DB
	var rdb *redis.Client
	var idxConnection *pinecone.IndexConnection
	var milvusClient client.Client
	var s3Client *s3.Client
	var firebaseApp *firebase.App

	conf.SetupConfig(ctx, &db, &rdb, &idxConnection, &milvusClient, &firebaseApp, &s3Client)

	boil.SetDB(db)
	//boil.DebugMode = true

	//// 차트 초기화
	//go handler.InitializeChart(db, rdb)
	//
	//// cronjob 추가
	//c := cron.New()
	//_, err = c.AddFunc("55 * * * *", func() {
	//	handler.ScheduleNextChart(db, rdb)
	//})
	//if err != nil {
	//	fmt.Println("Error scheduling task:", err)
	//	return
	//}
	//_, err = c.AddFunc("0 11 * * *", func() {
	//	handler.ScheduleNewSongs(db)
	//})
	//if err != nil {
	//	fmt.Println("Error scheduling task:", err)
	//	return
	//}
	//c.Start()

	r := router.SetupRouter(db, rdb, idxConnection, &milvusClient, firebaseApp, s3Client)

	// pprof를 위한 별도의 HTTP 서버 실행
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// 서버 실행

	srv := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: r.Handler(),
	}

	// 서버 실행이 블로킹(Blocking)되지 않도록 별도의 Go 루틴에서 실행하여 SIGTERM 감지를 위한 코드를 실행할 수 있도록 함.
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
	// make(chan os.Signal, 1) → OS에서 발생하는 신호(Signal)를 전달받는 채널 생성.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// <-quit → SIGTERM이 발생할 때까지 대기(Blocking).
	<-quit
	log.Println("Shutdown Server ...")

	//5초 동안 서버 종료를 기다릴 수 있는 컨텍스트 생성.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")

	//if err := r.Run(); err != nil {
	//	log.Fatalf("서버 실행 실패: %v", err)
	//}

}
