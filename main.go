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
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/pinecone-io/go-pinecone/pinecone"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
	"io"
	"log"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// main.go 상단 어딘가에 선언 (main 밖)
type otelLogWriter struct {
	logger *slog.Logger
}

func (w *otelLogWriter) Write(p []byte) (n int, err error) {
	w.logger.Info(string(p), "message", string(p))
	return len(p), nil
}

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
	// Instantiate a new slog logger
	ctx := context.Background()

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
			logrus.WithContext(ctx).Fatal("Failed to start profiler: ", err)
		}
		defer profiler.Stop()
	}

	otelShutdown, err := conf.SetupOTelSDK(ctx)
	if err != nil {
		return
	}
	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	logger := otelslog.NewLogger("Singsong")
	logger.Info("Hello from OpenTelemetry logs!", "orderID", 12345)

	var db *sql.DB
	var rdb *redis.Client
	var idxConnection *pinecone.IndexConnection
	var milvusClient client.Client
	var s3Client *s3.Client
	var firebaseApp *firebase.App

	conf.SetupConfig(ctx, &db, &rdb, &idxConnection, &milvusClient, &firebaseApp, &s3Client)

	boil.SetDB(db)
	boil.DebugMode = true

	// GIN 로그를 OTel 로그로 redirect
	gin.DefaultWriter = io.MultiWriter(&otelLogWriter{logger: logger}, os.Stdout, boil.DebugWriter)
	gin.DefaultErrorWriter = io.MultiWriter(&otelLogWriter{logger: logger}, os.Stderr, boil.DebugWriter)

	r := router.SetupRouter(db, rdb, idxConnection, &milvusClient, firebaseApp, s3Client)
	log.Println("router setup complete")

	pprofServer := &http.Server{
		Addr: "0.0.0.0:6060",
	}

	// pprof 서버 실행
	go func() {
		if err1 := pprofServer.ListenAndServe(); err1 != nil && err1 != http.ErrServerClosed {
			log.Fatalf("pprof listen: %s\n", err1)
		}
	}()

	// 메인 서버 실행
	srv := &http.Server{
		Addr:    ":8080",
		Handler: otelhttp.NewHandler(r.Handler(), "singsong-server"),
	}

	go func() {
		log.Println("server about to listen on :8080")
		if err2 := srv.ListenAndServe(); err2 != nil && err2 != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err2)
		}
	}()

	// Shutdown 처리
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	if err := pprofServer.Shutdown(ctx); err != nil {
		log.Fatal("pprof Shutdown:", err)
	}
}
