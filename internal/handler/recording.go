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

type RecordSongsRequest struct {
	Title       string `form:"title" json:"title"`
	Description string `form:"description" json:"description"`
	SongInfoId  int64  `form:"songInfoId" json:"songInfoId"`
	IsPublic    bool   `form:"isPublic" json:"isPublic"`
}

// RecordSong godoc
// @Summary      노래 녹음하기
// @Description  MP3파일을 받고 S3에 저장과 동시에 Presigned URL을 DB에 저장
// @Tags         Record
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "MP3 파일"
// @Param        title formData string true "노래 제목"
// @Param        description formData string true "설명"
// @Param        songInfoId formData int64 true "노래 정보 ID"
// @Param        isPublic formData bool true "공개 여부"
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
		file, handler, err := c.FormFile("file")
		if err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "파일 수신 실패", nil)
			return
		}
		defer file.Close()

		// 추가 메타데이터 수신
		var recordSongsRequest RecordSongsRequest
		if err := c.Bind(&recordSongsRequest); err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - "+err.Error(), nil)
			return
		}

		// S3 파일 이름 설정 및 업로드
		fileName := fmt.Sprintf("%s/%s", memberId, handler.Filename)
		_, err = s3Client.PutObject(c, &s3.PutObjectInput{
			Bucket: &S3BucketName,
			Key:    &fileName,
			Body:   file,
		})
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "S3 업로드 실패", nil)
			return
		}

		// S3 URL 생성
		s3URL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", S3BucketName, fileName)

		// URL을 DB에 저장
		_, err = db.Exec("INSERT INTO presigned_urls (url, title, description, song_info_id, is_public) VALUES (?, ?, ?, ?, ?)",
			s3URL, recordSongsRequest.Title, recordSongsRequest.Description, recordSongsRequest.SongInfoId, recordSongsRequest.IsPublic)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "DB 저장 실패", nil)
			return
		}

		// 성공 응답
		pkg.BaseResponse(c, http.StatusOK, "노래 녹음 성공", gin.H{
			"url": s3URL,
		})
	}
}
