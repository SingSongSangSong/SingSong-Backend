package handler

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/pinecone-io/go-pinecone/pinecone"
	"github.com/redis/go-redis/v9"
)

func GetRecommendation(db *sql.DB, redisClient *redis.Client, idxConnection *pinecone.IndexConnection) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
