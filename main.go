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
	"log"
	"net/http"
	_ "net/http/pprof"
)

// @title           싱송생송 API
// @version         1.0
// @BasePath  /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	ctx := context.Background()

	var db *sql.DB
	var rdb *redis.Client
	var idxConnection *pinecone.IndexConnection
	conf.SetupConfig(ctx, &db, &rdb, &idxConnection)
	// SQLBoiler의 디버그 모드 활성화
	//boil.DebugMode = true

	r := router.SetupRouter(db, rdb, idxConnection)

	// pprof를 위한 별도의 HTTP 서버 실행
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// 서버 실행
	if err := r.Run(); err != nil {
		log.Fatalf("서버 실행 실패: %v", err)
	}

}
