package handler

import (
	"SingSong-Server/internal/db/mysql"
	"SingSong-Server/internal/pkg"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"net/http"
	"strconv"
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
			qm.OrderBy("CASE WHEN song_name LIKE ? THEN 1 WHEN song_name LIKE ? THEN 2 ELSE 3 END", searchKeyword, searchKeyword+"%"),
			qm.Limit(10),
		).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// 아티스트 이름으로 검색
		songsWithArtist, err := mysql.SongInfos(
			qm.Where("artist_name LIKE ?", "%"+searchKeyword+"%"),
			qm.OrderBy("CASE WHEN artist_name LIKE ? THEN 1 WHEN artist_name LIKE ? THEN 2 ELSE 3 END", searchKeyword, searchKeyword+"%"),
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

var (
	defaultSearchSize = "20"
	defaultSearchPage = "1"
)

type songSearchPageResponse struct {
	Songs    []songSearchInfoResponse `json:"songs"`
	NextPage int                      `json:"nextPage"`
}

// SearchSongsByArtist godoc
// @Summary      가수로 노래 검색 API
// @Description  가수로 노래 검색 API, 아티스트 이름을 검색합니다. \n 검색 결과는 노래 제목, 아티스트 이름, 앨범명, 노래 번호를 반환합니다.
// @Tags         Search
// @Accept       json
// @Produce      json
// @Param        keyword query string true "검색 키워드"
// @Param        page query int false "현재 조회할 노래 목록의 쪽수. 입력하지 않는다면 기본값인 1쪽을 조회"
// @Param        size query int false "한번에 조회할 노래 개수. 입력하지 않는다면 기본값인 20개씩 조회"
// @Success      200 {object} pkg.BaseResponseStruct{data=songSearchPageResponse} "성공"
// @Failure      400 {object} pkg.BaseResponseStruct{data=nil} "실패"
// @Router       /search/artist-name [get]
func SearchSongsByArist(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 검색어를 쿼리 파라미터에서 가져오기
		searchKeyword := c.Query("keyword")
		if searchKeyword == "" {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - cannot find keyword in query", nil)
			return
		}
		pageValue := c.Query("page")
		if pageValue == "" {
			pageValue = defaultSearchPage
		}
		sizeValue := c.Query("size")
		if sizeValue == "" {
			sizeValue = defaultSearchSize
		}

		//page, size를 숫자로 변환
		page, err := strconv.Atoi(pageValue)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - cannot convert page to int", nil)
			return
		}
		size, err := strconv.Atoi(sizeValue)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - cannot convert size to int", nil)
			return
		}

		// 아티스트 이름으로 검색
		offset := (page - 1) * size

		songsWithArtist, err := mysql.SongInfos(
			qm.Where("artist_name LIKE ?", "%"+searchKeyword+"%"),
			qm.OrderBy("CASE WHEN artist_name LIKE ? THEN 1 WHEN artist_name LIKE ? THEN 2 ELSE 3 END", searchKeyword, searchKeyword+"%"),
			qm.Limit(size),
			qm.Offset(offset),
		).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// 응답 데이터 생성
		songs := make([]songSearchInfoResponse, 0, len(songsWithArtist))
		for _, song := range songsWithArtist {
			songs = append(songs, songSearchInfoResponse{
				SongNumber: song.SongNumber,
				SongName:   song.SongName,
				SingerName: song.ArtistName,
				SongInfoId: song.SongInfoID,
				Album:      song.Album.String,
				IsMr:       song.IsMR.Bool,
			})
		}
		response := songSearchPageResponse{
			Songs:    songs,
			NextPage: page + 1,
		}
		// 응답 반환
		pkg.BaseResponse(c, http.StatusOK, "ok", response)
	}
}

// SearchSongsBySongName godoc
// @Summary      노래 제목으로 노래 검색 API
// @Description  노래 제목으로 노래 검색 API, 노래 제목을 검색합니다. \n 검색 결과는 노래 제목, 아티스트 이름, 앨범명, 노래 번호를 반환합니다.
// @Tags         Search
// @Accept       json
// @Produce      json
// @Param        keyword query string true "검색 키워드"
// @Param        page query int false "현재 조회할 노래 목록의 쪽수. 입력하지 않는다면 기본값인 1쪽을 조회"
// @Param        size query int false "한번에 조회할 노래 개수. 입력하지 않는다면 기본값인 20개씩 조회"
// @Success      200 {object} pkg.BaseResponseStruct{data=songSearchPageResponse} "성공"
// @Failure      400 {object} pkg.BaseResponseStruct{data=nil} "실패"
// @Router       /search/song-name [get]
func SearchSongsBySongName(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 검색어를 쿼리 파라미터에서 가져오기
		searchKeyword := c.Query("keyword")
		if searchKeyword == "" {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - cannot find keyword in query", nil)
			return
		}
		pageValue := c.Query("page")
		if pageValue == "" {
			pageValue = defaultSearchPage
		}
		sizeValue := c.Query("size")
		if sizeValue == "" {
			sizeValue = defaultSearchSize
		}

		//page, size를 숫자로 변환
		page, err := strconv.Atoi(pageValue)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - cannot convert page to int", nil)
			return
		}
		size, err := strconv.Atoi(sizeValue)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - cannot convert size to int", nil)
			return
		}

		// 노래 이름으로 검색
		offset := (page - 1) * size
		songsWithName, err := mysql.SongInfos(
			qm.Where("song_name LIKE ?", "%"+searchKeyword+"%"),
			qm.OrderBy("CASE WHEN song_name LIKE ? THEN 1 WHEN song_name LIKE ? THEN 2 ELSE 3 END", searchKeyword, searchKeyword+"%"),
			qm.Limit(size),
			qm.Offset(offset),
		).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// 응답 데이터 생성
		songs := make([]songSearchInfoResponse, 0, len(songsWithName))
		for _, song := range songsWithName {
			songs = append(songs, songSearchInfoResponse{
				SongNumber: song.SongNumber,
				SongName:   song.SongName,
				SingerName: song.ArtistName,
				SongInfoId: song.SongInfoID,
				Album:      song.Album.String,
				IsMr:       song.IsMR.Bool,
			})
		}
		response := songSearchPageResponse{
			Songs:    songs,
			NextPage: page + 1,
		}
		// 응답 반환
		pkg.BaseResponse(c, http.StatusOK, "ok", response)
	}
}

// SearchSongsBySongNumber godoc
// @Summary      노래 번호로 노래 검색 API
// @Description  노래 번호로 노래 검색 API, 노래 번호를 검색합니다. \n 검색 결과는 노래 제목, 아티스트 이름, 앨범명, 노래 번호를 반환합니다.
// @Tags         Search
// @Accept       json
// @Produce      json
// @Param        keyword query string true "검색 키워드"
// @Param        page query int false "현재 조회할 노래 목록의 쪽수. 입력하지 않는다면 기본값인 1쪽을 조회. 현재는 노래 번호가 정확히 일치하는 1개만 반환하기 때문에 무의미"
// @Param        size query int false "한번에 조회할 노래 개수. 입력하지 않는다면 기본값인 20개씩 조회. 현재는 노래 번호가 정확히 일치하는 1개만 반환하기 때문에 무의미"
// @Success      200 {object} pkg.BaseResponseStruct{data=songSearchPageResponse} "성공"
// @Failure      400 {object} pkg.BaseResponseStruct{data=nil} "실패"
// @Router       /search/song-number [get]
func SearchSongsBySongNumber(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 검색어를 쿼리 파라미터에서 가져오기
		searchKeyword := c.Query("keyword")
		if searchKeyword == "" {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - cannot find keyword in query", nil)
			return
		}
		pageValue := c.Query("page")
		if pageValue == "" {
			pageValue = defaultSearchPage
		}

		//page를 숫자로 변환
		page, err := strconv.Atoi(pageValue)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - cannot convert page to int", nil)
			return
		}

		// 노래 번호로 검색
		songsWithNumber, err := mysql.SongInfos(
			qm.Where("song_number = ?", searchKeyword),
		).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// 응답 데이터 생성
		songs := make([]songSearchInfoResponse, 0, len(songsWithNumber))
		for _, song := range songsWithNumber {
			songs = append(songs, songSearchInfoResponse{
				SongNumber: song.SongNumber,
				SongName:   song.SongName,
				SingerName: song.ArtistName,
				SongInfoId: song.SongInfoID,
				Album:      song.Album.String,
				IsMr:       song.IsMR.Bool,
			})
		}
		response := songSearchPageResponse{
			Songs:    songs,
			NextPage: page + 1,
		}
		// 응답 반환
		pkg.BaseResponse(c, http.StatusOK, "ok", response)
	}
}
