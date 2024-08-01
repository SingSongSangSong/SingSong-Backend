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
	log "github.com/sirupsen/logrus"
)

// @title           싱송생송 API
// @version         1.0
// @BasePath  /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	ctx := context.Background()

	log.SetFormatter(&log.JSONFormatter{})

	var db *sql.DB
	var rdb *redis.Client
	var idxConnection *pinecone.IndexConnection
	conf.SetupConfig(ctx, &db, &rdb, &idxConnection)

	r := router.SetupRouter(db, rdb, idxConnection)

	log.WithFields(log.Fields{"string": "foo", "int": 1, "float": 1.1}).Info("My first event from golang to stdout")
	log.WithFields(log.Fields{"string": "231", "int": 2, "float": 2.1}).Error("My second event from golang to stdout")

	// 서버 실행
	if err := r.Run(); err != nil {
		log.Fatalf("서버 실행 실패: %v", err)
	}

}
