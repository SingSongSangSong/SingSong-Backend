package handler

import (
	"SingSong-Server/internal/db/mysql"
	"SingSong-Server/internal/pkg"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type newSongInfo struct {
	SongNumber        int    `json:"songNumber"`
	SongName          string `json:"songName"`
	SingerName        string `json:"singerName"`
	Album             string `json:"album"`
	IsKeep            bool   `json:"isKeep"`
	SongInfoId        int64  `json:"songId"`
	IsMr              bool   `json:"isMr"`
	IsLive            bool   `json:"isLive"`
	KeepCount         int    `json:"keepCount"`
	CommentCount      int    `json:"commentCount"`
	MelonLink         string `json:"melonLink"`
	IsRecentlyUpdated bool   `json:"isRecentlyUpdated"`
	LyricsYoutubeLink string `json:"lyricsYoutubeLink"`
	TJYoutubeLink     string `json:"tjYoutubeLink"`
	LyricsVideoID     string `json:"lyricsVideoId"`
	TJVideoID         string `json:"tjVideoId"`
}

type newSongInfoResponse struct {
	Songs      []newSongInfo `json:"songs"`
	LastCursor int64         `json:"lastCursor"`
}

// ListNewSongs godoc
// @Summary      최근 한달간의 신곡을 조회 (최신순 조회)
// @Description  최근 한달간의 신곡을 최신순으로 조회합니다. 최근 일주일동안 추가된 신곡은 isRecentlyUpdated = true 입니다.
// @Tags         New Song
// @Accept       json
// @Produce      json
// @Param        cursor query int false "마지막에 조회했던 커서의 songId(이전 요청에서 lastCursor값을 주면 됨), 없다면 default로 가장 최신곡부터 조회"
// @Param        size query int false "한번에 가져욜 노래 개수. 입력하지 않는다면 기본값인 20개씩 조회"
// @Success      200 {object} pkg.BaseResponseStruct{data=newSongInfoResponse} "성공"
// @Failure      400 "query param 값이 들어왔는데, 숫자가 아니라면 400 실패"
// @Failure      500 "서버 에러일 경우 500 실패"
// @Router       /v1/songs/new [get]
// @Security BearerAuth
func ListNewSongs(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		memberId, exists := c.Get("memberId")
		if !exists {
			pkg.SendToSentryWithStack(c, fmt.Errorf("memberId not found in context"))
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - memberId not found", nil)
			return
		}

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

		// 페이징 처리된 신곡 가져오기
		// 1. 페이징 처리된 신곡 가져오기 (메인 쿼리)
		songInfos, err := mysql.SongInfos(
			qm.Where("song_info_id < ?", cursorInt),
			qm.And("created_at >= DATE_SUB(NOW(), INTERVAL 1 MONTH)"),
			qm.OrderBy("song_info_id DESC"),
			qm.Limit(sizeInt),
		).All(c, db)
		if err != nil {
			pkg.SendToSentryWithStack(c, err)
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		if len(songInfos) == 0 {
			pkg.BaseResponse(c, http.StatusOK, "ok", []interface{}{})
			return
		}

		// 2. song_info_id 리스트 생성
		songInfoIDs := make([]interface{}, len(songInfos))
		for i, v := range songInfos {
			songInfoIDs[i] = v.SongInfoID
		}

		// 3. Goroutine을 사용하여 Keep, Comment, Keep Count 조회 병렬 실행
		var wg sync.WaitGroup
		wg.Add(3) // 3개의 고루틴 실행

		// 에러 처리 및 결과 저장용 채널
		errChan := make(chan error, 3)
		keepMap := make(map[int64]bool)
		commentCountMap := make(map[int64]int)
		keepCountMap := make(map[int64]int)

		// Keep 여부 조회
		go func() {
			defer wg.Done()
			keepLists, err := mysql.KeepLists(qm.Where("member_id = ?", memberId)).All(c, db)
			if err != nil {
				errChan <- err
				return
			}

			keepListIDs := make([]interface{}, len(keepLists))
			for i, v := range keepLists {
				keepListIDs[i] = v.KeepListID
			}

			keepSongs, err := mysql.KeepSongs(
				qm.WhereIn("keep_list_id = ?", keepListIDs...),
				qm.AndIn("song_info_id IN ?", songInfoIDs...),
			).All(c, db)
			if err != nil {
				errChan <- err
				return
			}

			for _, keepSong := range keepSongs {
				keepMap[keepSong.SongInfoID] = true
			}
		}()

		// 댓글 수 조회
		go func() {
			defer wg.Done()
			commentsCounts, err := mysql.Comments(
				qm.WhereIn("song_info_id IN ?", songInfoIDs...),
				qm.And("deleted_at IS NULL"),
			).All(c, db)
			if err != nil {
				errChan <- err
				return
			}

			for _, comment := range commentsCounts {
				commentCountMap[comment.SongInfoID]++
			}
		}()

		// Keep 개수 조회
		go func() {
			defer wg.Done()
			keepCounts, err := mysql.KeepSongs(
				qm.WhereIn("song_info_id IN ?", songInfoIDs...),
				qm.And("deleted_at IS NULL"),
			).All(c, db)
			if err != nil {
				errChan <- err
				return
			}

			for _, keep := range keepCounts {
				keepCountMap[keep.SongInfoID]++
			}
		}()

		// 모든 고루틴이 종료될 때까지 대기
		wg.Wait()
		close(errChan)

		// 에러가 있으면 반환
		for err := range errChan {
			if err != nil {
				pkg.SendToSentryWithStack(c, err)
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}
		}

		newSongs := make([]newSongInfo, 0, sizeInt)
		sevenDaysAgo := time.Now().AddDate(0, 0, -7)

		for _, song := range songInfos {
			newSongs = append(newSongs, newSongInfo{
				SongNumber:        song.SongNumber,
				SongName:          song.SongName,
				SingerName:        song.ArtistName,
				Album:             song.Album.String,
				IsKeep:            keepMap[song.SongInfoID],
				SongInfoId:        song.SongInfoID,
				IsMr:              song.IsMR.Bool,
				IsLive:            song.IsLive.Bool,
				KeepCount:         keepCountMap[song.SongInfoID],
				CommentCount:      commentCountMap[song.SongInfoID],
				MelonLink:         CreateMelonLinkByMelonSongId(song.MelonSongID),
				IsRecentlyUpdated: song.CreatedAt.Time.After(sevenDaysAgo), //song.CreatedAt이 7일 이내인지
				LyricsYoutubeLink: song.LyricsVideoLink.String,
				TJYoutubeLink:     song.TJYoutubeLink.String,
				LyricsVideoID:     ExtractVideoID(song.LyricsVideoLink.String),
				TJVideoID:         ExtractVideoID(song.TJYoutubeLink.String),
			})
		}

		// 다음 페이지를 위한 커서 값 설정
		var lastCursor int64 = 0
		if len(newSongs) > 0 {
			lastCursor = newSongs[len(newSongs)-1].SongInfoId
		}

		response := newSongInfoResponse{
			Songs:      newSongs,
			LastCursor: lastCursor,
		}

		// 응답 반환
		pkg.BaseResponse(c, http.StatusOK, "ok", response)
	}
}
