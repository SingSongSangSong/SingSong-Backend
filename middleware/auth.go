package middleware

import (
	"SingSong-Server/conf"
	"SingSong-Server/internal/db/mysql"
	"SingSong-Server/internal/pkg"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"net/http"
	"strings"
)

var (
	secretKey = []byte(conf.AuthConfigInstance.SECRET_KEY)
)

type claims struct {
	Email    string
	Provider string
	jwt.StandardClaims
}

func AuthMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			pkg.BaseResponse(c, http.StatusUnauthorized, "error - Authorization header is required", nil)
			c.Abort() // 다음 처리를 중단
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			pkg.BaseResponse(c, http.StatusUnauthorized, "error - invalid Authorization header format", nil)
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("invalid signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		})

		if err != nil {
			pkg.BaseResponse(c, http.StatusUnauthorized, "error - invalid token "+err.Error(), nil)
			c.Abort()
			return
		}

		// 컨텍스트에 email과 provider를 저장
		var email, provider string
		if claims, ok := token.Claims.(*claims); ok && token.Valid {
			email = claims.Email
			provider = claims.Provider
			c.Set("email", claims.Email)
			c.Set("provider", claims.Provider)
		} else {
			pkg.BaseResponse(c, http.StatusUnauthorized, "error - invalid token", nil)
			c.Abort()
			return
		}

		one, err := mysql.Members(qm.Where("email = ? AND provider = ?", email, provider)).One(c, db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusUnauthorized, "error - "+err.Error(), nil)
			c.Abort()
			return
		}
		c.Set("memberId", one.ID)
		c.Next() // 다음 핸들러로 고고싱
	}
}
