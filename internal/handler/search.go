package handler

import (
	"SingSong-Server/internal/db/mysql"
	"SingSong-Server/internal/pkg"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"net/http"
)

type songSearchInfoResponse struct {
	SongNumber int    `json:"songNumber"`
	SongName   string `json:"songName"`
	SingerName string `json:"singerName"`
	SongInfoId int64  `json:"songId"`
	Album      string `json:"album"`
	IsMr       bool   `json:"isMr"`
}

type songSearchInfoResponses struct {
	SongName   []songSearchInfoResponse `json:"songName"`
	ArtistName []songSearchInfoResponse `json:"artistName"`
	SongNumber []songSearchInfoResponse `json:"songNumber"`
}

// SearchSongs godoc
// @Summary      노래 검색 API
// @Description  노래 검색 API, 노래 제목 또는 아티스트 이름을 검색합니다. \n 검색 결과는 노래 제목, 아티스트 이름, 앨범명, 노래 번호를 반환합니다.
// @Tags         Search
// @Accept       json
// @Produce      json
// @Param        searchKeyword path string true "검색 키워드"
// @Success      200 {object} pkg.BaseResponseStruct{data=songSearchInfoResponses} "성공"
// @Failure      400 {object} pkg.BaseResponseStruct{data=nil} "실패 - 빈 리스트 반환"
// @Router       /search/{searchKeyword} [get]
func SearchSongs(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 검색어를 URL 파라미터에서 가져오기
		searchKeyword := c.Param("searchKeyword")

		// 노래 이름으로 검색
		songsWithName, err := mysql.SongInfos(
			qm.Where("song_name LIKE ?", "%"+searchKeyword+"%"),
			qm.Limit(10),
		).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// 아티스트 이름으로 검색
		songsWithArtist, err := mysql.SongInfos(
			qm.Where("artist_name LIKE ?", "%"+searchKeyword+"%"),
			qm.Limit(10),
		).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// 노래 번호로 검색
		songsWithNumber, err := mysql.SongInfos(
			qm.Where("song_number = ?", searchKeyword),
			qm.Limit(10),
		).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// 응답 데이터 생성
		response := songSearchInfoResponses{
			SongName:   make([]songSearchInfoResponse, 0),
			ArtistName: make([]songSearchInfoResponse, 0),
			SongNumber: make([]songSearchInfoResponse, 0),
		}

		// 노래 이름으로 검색한 결과를 response 추가
		for _, song := range songsWithName {
			response.SongName = append(response.SongName, songSearchInfoResponse{
				SongNumber: song.SongNumber,
				SongName:   song.SongName,
				SingerName: song.ArtistName,
				SongInfoId: song.SongInfoID,
				Album:      song.Album.String,
				IsMr:       song.IsMR.Bool,
			})
		}

		// 아티스트 이름으로 검색한 결과를 response 추가
		for _, song := range songsWithArtist {
			response.ArtistName = append(response.ArtistName, songSearchInfoResponse{
				SongNumber: song.SongNumber,
				SongName:   song.SongName,
				SingerName: song.ArtistName,
				SongInfoId: song.SongInfoID,
				Album:      song.Album.String,
				IsMr:       song.IsMR.Bool,
			})
		}

		// 노래 번호로 검색한 결과를 response에 추가
		for _, song := range songsWithNumber {
			response.SongNumber = append(response.SongNumber, songSearchInfoResponse{
				SongNumber: song.SongNumber,
				SongName:   song.SongName,
				SingerName: song.ArtistName,
				SongInfoId: song.SongInfoID,
				Album:      song.Album.String,
				IsMr:       song.IsMR.Bool,
			})
		}

		// 응답 반환
		pkg.BaseResponse(c, http.StatusOK, "ok", response)
	}
}
