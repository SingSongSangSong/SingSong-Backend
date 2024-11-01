package handler

import (
	"SingSong-Server/conf"
	"SingSong-Server/internal/pkg"
	"database/sql"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	S3BucketName = conf.AWSConfigInstance.S3BucketName
)

type result struct {
	Output *s3.PutObjectOutput
	Err    error
}

// RecordSong godoc
// @Summary      노래 녹음하기
// @Description  버튼을 클릭할때부터 MP3파일을 받고 S3에 저장과 동시에 Presigned URL DB에 저장
// @Tags         Record
// @Accept       mp3
// @Produce      json
// @Param        file formData file true "MP3 파일"
// @Success      200 {object} pkg.BaseResponseStruct{data=UserProfileResponse} "성공"
// @Router       /v1/record/song [post]
// @Security BearerAuth
func RecordSong(db *sql.DB, s3Client *s3.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		memberId, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - memberId not found", nil)
			return
		}

		// 파일 수신
		file, handler, err := c.Request.FormFile("file")
		if err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "파일 수신 실패", nil)
			return
		}
		defer file.Close()

		// 파일 이름 설정 및 S3 업로드
		fileName := fmt.Sprintf("%s/%s", memberId, handler.Filename)

		results := make(chan result, 2)
		go func() {
			output, err := s3Client.PutObject(c, &s3.PutObjectInput{
				Bucket: &S3BucketName,
				Key:    &fileName,
				Body:   file,
			})
			results <- result{Output: output, Err: err}
		}()

		// Presigned URL을 DB에 저장
		_, err = db.Exec("INSERT INTO presigned_urls (url) VALUES (?)", uploadURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB 저장 실패"})
			return
		}

		pkg.BaseResponse(c, http.StatusOK, "노래 녹음 성공", nil)
	}
}
