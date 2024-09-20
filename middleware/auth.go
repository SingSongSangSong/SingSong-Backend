package middleware

import (
	"SingSong-Server/conf"
	"SingSong-Server/internal/handler"
	"SingSong-Server/internal/pkg"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"strings"
)

var (
	secretKey = []byte(conf.AuthConfigInstance.SECRET_KEY)
)

func AuthMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			pkg.BaseResponse(c, http.StatusUnauthorized, "error - Authorization header is required", nil)
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			pkg.BaseResponse(c, http.StatusUnauthorized, "error - invalid Authorization header format", nil)
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &handler.Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("error - invalid signing method")
			}
			return secretKey, nil
		})

		if err != nil {
			pkg.BaseResponse(c, http.StatusUnauthorized, "error - invalid token", nil)
			c.Abort()
			return
		}

		var gender, birthYear string
		var memberId int64

		if claims, ok := token.Claims.(*handler.Claims); ok && token.Valid {
			memberId = claims.MemberId
			gender = claims.Gender
			birthYear = claims.BirthYear

		} else {
			pkg.BaseResponse(c, http.StatusUnauthorized, "error - invalid token", nil)
			c.Abort()
			return
		}

		log.Printf("memberId: " + fmt.Sprintf("%d", memberId))

		c.Set("memberId", memberId)
		c.Set("gender", gender)
		c.Set("birthYear", birthYear)
		c.Next()
	}
}
