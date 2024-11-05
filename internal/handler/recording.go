package handler

import (
	"SingSong-Server/conf"
	"SingSong-Server/internal/db/mysql"
	"SingSong-Server/internal/pkg"
	"context"
	"database/sql"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"net/http"
	"strconv"
	"time"
)

var (
	S3BucketName = conf.AWSConfigInstance.S3BucketName
)

type RecordSongsRequest struct {
	Title    string `form:"title" json:"title"`
	SongId   int64  `form:"songId" json:"songId"`
	IsPublic bool   `form:"isPublic" json:"isPublic"`
}

// RecordSong godoc
// @Summary      노래 녹음하기
// @Description  MP3파일을 받고 S3에 저장과 동시에 URL을 DB에 저장
// @Tags         Record
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "MP3 파일"
// @Param        title formData string true "노래 제목"
// @Param        songId formData int64 true "노래 정보 ID"
// @Param        isPublic formData bool true "공개 여부"
// @Success      200 {object} pkg.BaseResponseStruct{data=nil} "성공"
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
		getFile, err := c.FormFile("file")
		if err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "파일 수신 실패", nil)
			return
		}

		// 파일 형식을 확인해야한다. 파일의 형식이 오직 Mp3만 허용한다.
		if getFile.Header.Get("Content-Type") != "audio/mpeg" {
			pkg.BaseResponse(c, http.StatusBadRequest, "파일 형식이 올바르지 않습니다. MP3 파일만 허용됩니다.", nil)
			return
		}

		// 파일 열기
		file, err := getFile.Open()
		if err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "파일 열기 실패", nil)
			return
		}

		// 추가 메타데이터 수신
		var recordSongsRequest RecordSongsRequest
		if err := c.Bind(&recordSongsRequest); err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - "+err.Error(), nil)
			return
		}

		// 현재 시간 가져오기
		currentTime := time.Now().Format("20060102_150405") // YYYYMMDD_HHMMSS 형식

		fileName := ""

		if conf.Env == conf.LocalMode || conf.Env == "" {
			fileName = fmt.Sprintf("local/%d/%d/%s.mp3", memberId.(int64), recordSongsRequest.SongId, currentTime)
		} else if conf.Env == conf.TestMode {
			fileName = fmt.Sprintf("test/%d/%d/%s.mp3", memberId.(int64), recordSongsRequest.SongId, currentTime)
		} else if conf.Env == conf.ProductionMode {
			fileName = fmt.Sprintf("prod/%d/%d/%s.mp3", memberId.(int64), recordSongsRequest.SongId, currentTime)
		} else {
			pkg.BaseResponse(c, http.StatusInternalServerError, "환경 변수가 설정되지 않았습니다.", nil)
			return
		}
		// S3 URL 생성
		s3URL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", S3BucketName, fileName)

		// S3 파일 이름 설정 및 업로드
		_, err = s3Client.PutObject(c, &s3.PutObjectInput{
			Bucket: &S3BucketName,
			Key:    &fileName,
			Body:   file,
		})
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "S3 업로드 실패", nil)
			return
		}

		// 노래 정보 DB에 저장
		songRecording := mysql.SongRecording{MemberID: memberId.(int64), Title: recordSongsRequest.Title, SongInfoID: recordSongsRequest.SongId, IsPublic: null.BoolFrom(recordSongsRequest.IsPublic), RecordingLink: s3URL}
		err = songRecording.Insert(c, db, boil.Infer())
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "DB 저장 실패", nil)
			return
		}

		// 성공 응답
		pkg.BaseResponse(c, http.StatusOK, "노래 녹음 성공", nil)
	}
}

type SongRecording struct {
	SongRecordingID   int64  `json:"songRecordingId"`
	Title             string `json:"title"`
	IsPublic          bool   `json:"isPublic"`
	SongId            int64  `json:"SongId"`
	SongNumber        int    `json:"songNumber"`
	SongName          string `json:"songName"`
	SingerName        string `json:"singerName"`
	Album             string `json:"album"`
	IsMr              bool   `json:"isMr"`
	IsLive            bool   `json:"isLive"`
	MelonLink         string `json:"melonLink"`
	LyricsYoutubeLink string `json:"lyricsYoutubeLink"`
	TJYoutubeLink     string `json:"tjYoutubeLink"`
	CreatedAt         string `json:"createdAt"`
}

type GetRecordingsResponse struct {
	SongRecordings []SongRecording `json:"songRecordings"`
	LastCursor     int64           `json:"lastCursor"`
}

// GetMyRecordings godoc
// @Summary      내 녹음 목록 조회
// @Description  내가 녹음한 노래 목록을 조회한다
// @Tags         Record
// @Accept       json
// @Produce      json
// @Param        cursor query int false "마지막에 조회했던 커서의 SongRecordingId(이전 요청에서 lastCursor값을 주면 됨), 없다면 default로 가장 최신곡부터 조회"
// @Param        size query int false "한번에 가져욜 노래 개수. 입력하지 않는다면 기본값인 20개씩 조회"
// @Success      200 {object} pkg.BaseResponseStruct{data=GetRecordingsResponse} "성공"
// @Router       /v1/record/list [get]
// @Security BearerAuth
func GetMyRecordings(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		memberId, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - memberId not found", nil)
			return
		}

		// cursor, size 파라미터 수신
		sizeStr := c.DefaultQuery("size", defaultSize)
		sizeInt, err := strconv.Atoi(sizeStr)
		if err != nil || sizeInt < 0 {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid size parameter", nil)
			return
		}

		cursorStr := c.DefaultQuery("cursor", "9223372036854775807") //int64 최대값
		cursorInt, err := strconv.Atoi(cursorStr)
		if err != nil || cursorInt <= 0 {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid cursor parameter", nil)
			return
		}

		recordings, err := mysql.SongRecordings(
			qm.Load(mysql.SongRecordingRels.SongInfo),
			qm.Where("member_id = ? AND song_recording_id < ? AND deleted_at IS NULL", memberId, cursorInt),
			qm.OrderBy("created_at DESC"),
			qm.Limit(sizeInt),
		).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "DB 조회 실패", nil)
			return
		}

		// 녹음 목록이 없을 경우
		if len(recordings) == 0 {
			pkg.BaseResponse(c, http.StatusOK, "녹음 목록이 없습니다.", GetRecordingsResponse{[]SongRecording{}, 0})
			return
		}

		// 녹음 목록이 있을 경우
		songRecordings := make([]SongRecording, len(recordings))
		for i, recording := range recordings {
			songRecordings[i] = SongRecording{
				SongRecordingID:   recording.SongRecordingID,
				Title:             recording.Title,
				IsPublic:          recording.IsPublic.Bool,
				SongId:            recording.SongInfoID,
				SongNumber:        recording.R.SongInfo.SongNumber,
				SongName:          recording.R.SongInfo.SongName,
				SingerName:        recording.R.SongInfo.ArtistName,
				Album:             recording.R.SongInfo.Album.String,
				IsMr:              recording.R.SongInfo.IsMR.Bool,
				IsLive:            recording.R.SongInfo.IsLive.Bool,
				MelonLink:         CreateMelonLinkByMelonSongId(recording.R.SongInfo.MelonSongID),
				LyricsYoutubeLink: recording.R.SongInfo.LyricsVideoLink.String,
				TJYoutubeLink:     recording.R.SongInfo.TJYoutubeLink.String,
				CreatedAt:         recording.CreatedAt.Time.String(),
			}
		}

		response := GetRecordingsResponse{
			SongRecordings: songRecordings,
			LastCursor:     songRecordings[len(songRecordings)-1].SongRecordingID,
		}

		pkg.BaseResponse(c, http.StatusOK, "녹음 목록 조회 성공", response)
	}
}

type GetDetailRecordingResponse struct {
	SongRecording SongRecording `json:"songRecording"`
	PreSignedURL  string        `json:"preSignedURL"`
}

// GetDetailRecording godoc
// @Summary      녹음 상세 조회
// @Description  녹음 상세 정보를 조회한다. 이때 바로 플레이가 가능하게끔 15분동안 접근이 가능한 Presigned URL을 생성하여 반환한다.이는 15분뒤에는 접근이 불가능해지는 링크이다.
// @Tags         Record
// @Accept       json
// @Produce      json
// @Param        songRecordingId path int true "녹음 ID"
// @Success      200 {object} pkg.BaseResponseStruct{data=GetDetailRecordingResponse} "성공"
// @Failure      500 "녹음 정보가 없는 경우" "error"
// @Router       /v1/record/{songRecordingId}/my [get]
// @Security BearerAuth
func GetDetailRecording(db *sql.DB, s3Client *s3.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// memberId 가져오기
		memberId, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - memberId not found", nil)
			return
		}

		// songRecordingId 가져오기
		songRecordingId := c.Param("songRecordingId")
		if songRecordingId == "" {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - songRecordingId not found", nil)
			return
		}

		// DB에서 녹음 정보 조회
		songRecording := mysql.SongRecordings(
			qm.Load(mysql.SongRecordingRels.SongInfo),
			qm.Where("song_recording_id = ? AND member_id = ? AND deleted_at IS NULL", songRecordingId, memberId),
		)
		songRecordings, err := songRecording.One(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "DB 조회 실패", nil)
			return
		}

		// 녹음 정보가 없는 경우 처리
		if songRecordings == nil {
			pkg.BaseResponse(c, http.StatusNotFound, "녹음 정보가 없습니다.", nil)
			return
		}

		// Presigned URL 생성
		bucketName := S3BucketName // 실제 버킷 이름
		// RecordingLink에서 Key 추출
		prefix := fmt.Sprintf("https://%s.s3.amazonaws.com/", S3BucketName)
		key := songRecordings.RecordingLink[len(prefix):] // URL Key 부분만 추출
		expiration := 15 * time.Minute                    // URL 만료 시간 설정

		preSignedURL, err := generatePresignedURL(s3Client, bucketName, key, expiration)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "Presigned URL 생성 실패", nil)
			return
		}

		songRecordingResponse := SongRecording{
			SongRecordingID:   songRecordings.SongRecordingID,
			Title:             songRecordings.Title,
			IsPublic:          songRecordings.IsPublic.Bool,
			SongId:            songRecordings.SongInfoID,
			SongNumber:        songRecordings.R.SongInfo.SongNumber,
			SongName:          songRecordings.R.SongInfo.SongName,
			SingerName:        songRecordings.R.SongInfo.ArtistName,
			Album:             songRecordings.R.SongInfo.Album.String,
			IsMr:              songRecordings.R.SongInfo.IsMR.Bool,
			IsLive:            songRecordings.R.SongInfo.IsLive.Bool,
			MelonLink:         CreateMelonLinkByMelonSongId(songRecordings.R.SongInfo.MelonSongID),
			LyricsYoutubeLink: songRecordings.R.SongInfo.LyricsVideoLink.String,
			TJYoutubeLink:     songRecordings.R.SongInfo.TJYoutubeLink.String,
			CreatedAt:         songRecordings.CreatedAt.Time.String(),
		}

		// 응답 데이터 생성
		response := GetDetailRecordingResponse{
			SongRecording: songRecordingResponse,
			PreSignedURL:  preSignedURL,
		}

		// 성공 응답
		pkg.BaseResponse(c, http.StatusOK, "녹음 상세 조회 성공", response)
	}
}

// DeleteMyRecording godoc
// @Summary      녹음 삭제
// @Description  녹음을 삭제 한다 하지만 AWS 보안상 삭제 권한을 주지 않아서 S3에서 삭제 하지 않는다.
// @Tags         Record
// @Accept       json
// @Produce      json
// @Param        songRecordingId path int true "녹음 ID"
// @Success      200 {object} pkg.BaseResponseStruct{data=nil} "성공"
// @Router       /v1/record/{songRecordingId}/my [delete]
// @Security BearerAuth
func DeleteMyRecording(db *sql.DB, s3Client *s3.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		songRecordingId := c.Param("songRecordingId")
		if songRecordingId == "" {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - songRecordingId not found", nil)
			return
		}

		// memberId 가져오기
		memberId, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - memberId not found", nil)
			return
		}

		// DB에서 녹음 정보 조회
		songRecording, err := mysql.SongRecordings(
			qm.Where("song_recording_id = ? AND member_id = ? AND deleted_at IS NULL", songRecordingId, memberId),
		).One(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "DB 조회 실패", nil)
			return
		}

		// deleted_at 필드 업데이트
		songRecording.DeletedAt = null.TimeFrom(time.Now()) // null.Time 사용 시 null.TimeFrom 사용
		_, err = songRecording.Update(c.Request.Context(), db, boil.Infer())
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "DB 삭제 실패", nil)
			return
		}

		//// S3에서 파일 삭제 기능을 넣을려고 하였으나 사용자에게 AWS 삭제권한을 주는것은 위험하다고 판단하여 주지 않음
		//prefix := fmt.Sprintf("https://%s.s3.amazonaws.com/", S3BucketName)
		//key := songRecording.RecordingLink[len(prefix):] // URL Key 부분만 추출
		//_, err = s3Client.DeleteObject(c, &s3.DeleteObjectInput{
		//	Bucket: &S3BucketName,
		//	Key:    &key,
		//})
		//if err != nil {
		//	pkg.BaseResponse(c, http.StatusInternalServerError, "S3 파일 삭제 실패", nil)
		//	return
		//}

		// 성공 응답
		pkg.BaseResponse(c, http.StatusOK, "녹음 삭제 성공", nil)
	}
}

// Presigned URL 생성 함수
func generatePresignedURL(s3Client *s3.Client, bucketName, key string, expiration time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(s3Client)

	// Presigned URL을 생성하는 요청
	req, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expiration))

	if err != nil {
		return "", fmt.Errorf("Presigned URL 생성 실패: %w", err)
	}

	return req.URL, nil
}
