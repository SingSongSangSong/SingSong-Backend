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

type SongSearchInfoV2Response struct {
	SongNumber        int    `json:"songNumber"`
	SongName          string `json:"songName"`
	SingerName        string `json:"singerName"`
	SongInfoId        int64  `json:"songId"`
	Album             string `json:"album"`
	IsMr              bool   `json:"isMr"`
	IsLive            bool   `json:"isLive"`
	MelonLink         string `json:"melonLink"`
	IsKeep            bool   `json:"isKeep"`
	LyricsYoutubeLink string `json:"lyricsYoutubeLink"`
	TJYoutubeLink     string `json:"tjYoutubeLink"`
}

type SongSearchInfoV2Responses struct {
	SongName   []SongSearchInfoV2Response `json:"songName"`
	ArtistName []SongSearchInfoV2Response `json:"artistName"`
	SongNumber []SongSearchInfoV2Response `json:"songNumber"`
}

// SearchSongsV2 godoc
// @Summary      노래 검색 API V2
// @Description  노래 검색 API V2, 노래 제목 또는 아티스트 이름을 검색합니다. \n 검색 결과는 노래 제목, 아티스트 이름, 앨범명, 노래 번호 및 추가적으로 IsKeep여부를 반환합니다.
// @Tags         Search
// @Accept       json
// @Produce      json
// @Param        searchKeyword path string true "검색 키워드"
// @Success      200 {object} pkg.BaseResponseStruct{data=SongSearchInfoV2Responses} "성공"
// @Failure      400 {object} pkg.BaseResponseStruct{data=nil} "실패 - 빈 리스트 반환"
// @Router       /v2/search/{searchKeyword} [get]
// @Security BearerAuth
func SearchSongsV2(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		memberId, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - memberId not found", nil)
			return
		}

		// 검색어를 URL 파라미터에서 가져오기
		searchKeyword := c.Param("searchKeyword")

		// 노래 이름으로 검색
		songsWithName, err := mysql.SongInfos(
			qm.Where("song_name LIKE ?", "%"+searchKeyword+"%"),
			qm.OrderBy("CASE WHEN song_name LIKE ? THEN 1 WHEN song_name LIKE ? THEN 2 ELSE 3 END, song_info_id", searchKeyword, searchKeyword+"%"),
			qm.Limit(10),
		).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// 아티스트 이름으로 검색
		songsWithArtist, err := mysql.SongInfos(
			qm.Where("artist_name LIKE ?", "%"+searchKeyword+"%"),
			qm.OrderBy("CASE WHEN artist_name LIKE ? THEN 1 WHEN artist_name LIKE ? THEN 2 ELSE 3 END, song_info_id", searchKeyword, searchKeyword+"%"),
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

		keepList, err := mysql.KeepLists(
			qm.Where("member_id = ?", memberId),
		).One(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		keepSongs, err := mysql.KeepSongs(
			qm.Where("keep_list_id = ?", keepList.KeepListID),
		).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		keepSongMap := make(map[int64]bool)
		for _, keepSong := range keepSongs {
			keepSongMap[keepSong.SongInfoID] = true
		}

		// 응답 데이터 생성
		response := SongSearchInfoV2Responses{
			SongName:   make([]SongSearchInfoV2Response, 0),
			ArtistName: make([]SongSearchInfoV2Response, 0),
			SongNumber: make([]SongSearchInfoV2Response, 0),
		}

		// 노래 이름으로 검색한 결과를 response 추가
		for _, song := range songsWithName {
			response.SongName = append(response.SongName, SongSearchInfoV2Response{
				SongNumber:        song.SongNumber,
				SongName:          song.SongName,
				SingerName:        song.ArtistName,
				SongInfoId:        song.SongInfoID,
				Album:             song.Album.String,
				IsMr:              song.IsMR.Bool,
				IsLive:            song.IsLive.Bool,
				MelonLink:         CreateMelonLinkByMelonSongId(song.MelonSongID),
				IsKeep:            keepSongMap[song.SongInfoID],
				LyricsYoutubeLink: song.LyricsVideoLink.String,
				TJYoutubeLink:     song.TJYoutubeLink.String,
			})
		}

		// 아티스트 이름으로 검색한 결과를 response 추가
		for _, song := range songsWithArtist {
			response.ArtistName = append(response.ArtistName, SongSearchInfoV2Response{
				SongNumber:        song.SongNumber,
				SongName:          song.SongName,
				SingerName:        song.ArtistName,
				SongInfoId:        song.SongInfoID,
				Album:             song.Album.String,
				IsMr:              song.IsMR.Bool,
				IsLive:            song.IsLive.Bool,
				MelonLink:         CreateMelonLinkByMelonSongId(song.MelonSongID),
				IsKeep:            keepSongMap[song.SongInfoID],
				LyricsYoutubeLink: song.LyricsVideoLink.String,
				TJYoutubeLink:     song.TJYoutubeLink.String,
			})
		}

		// 노래 번호로 검색한 결과를 response에 추가
		for _, song := range songsWithNumber {
			response.SongNumber = append(response.SongNumber, SongSearchInfoV2Response{
				SongNumber:        song.SongNumber,
				SongName:          song.SongName,
				SingerName:        song.ArtistName,
				SongInfoId:        song.SongInfoID,
				Album:             song.Album.String,
				IsMr:              song.IsMR.Bool,
				IsLive:            song.IsLive.Bool,
				MelonLink:         CreateMelonLinkByMelonSongId(song.MelonSongID),
				IsKeep:            keepSongMap[song.SongInfoID],
				LyricsYoutubeLink: song.LyricsVideoLink.String,
				TJYoutubeLink:     song.TJYoutubeLink.String,
			})
		}

		// 응답 반환
		pkg.BaseResponse(c, http.StatusOK, "ok", response)
	}
}

type SongSearchPageV2Response struct {
	Songs    []SongSearchInfoV2Response `json:"songs"`
	NextPage int                        `json:"nextPage"`
}

// SearchSongsByAristV2 godoc
// @Summary      가수로 노래 검색 API V2
// @Description  가수로 노래 검색 API, 아티스트 이름을 검색합니다. \n 검색 결과는 노래 제목, 아티스트 이름, 앨범명, 노래 번호 및 IsKeep여부를 반환합니다.
// @Tags         Search
// @Accept       json
// @Produce      json
// @Param        keyword query string true "검색 키워드"
// @Param        page query int false "현재 조회할 노래 목록의 쪽수. 입력하지 않는다면 기본값인 1쪽을 조회"
// @Param        size query int false "한번에 조회할 노래 개수. 입력하지 않는다면 기본값인 20개씩 조회"
// @Success      200 {object} pkg.BaseResponseStruct{data=SongSearchPageV2Response} "성공"
// @Failure      400 {object} pkg.BaseResponseStruct{data=nil} "실패"
// @Router       /v2/search/artist-name [get]
// @Security BearerAuth
func SearchSongsByAristV2(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		memberId, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - memberId not found", nil)
			return
		}

		// 검색어를 쿼리 파라미터에서 가져오기
		searchKeyword := c.Query("keyword")
		if searchKeyword == "" {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - cannot find keyword in query", nil)
			return
		}
		pageValue := c.DefaultQuery("page", defaultSearchPage)
		sizeValue := c.DefaultQuery("size", defaultSearchSize)

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
			qm.OrderBy("CASE WHEN artist_name LIKE ? THEN 1 WHEN artist_name LIKE ? THEN 2 ELSE 3 END, song_info_id", searchKeyword, searchKeyword+"%"),
			qm.Limit(size),
			qm.Offset(offset),
		).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		keepList, err := mysql.KeepLists(
			qm.Where("member_id = ?", memberId),
		).One(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		keepSongs, err := mysql.KeepSongs(
			qm.Where("keep_list_id = ?", keepList.KeepListID),
		).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		keepSongMap := make(map[int64]bool)
		for _, keepSong := range keepSongs {
			keepSongMap[keepSong.SongInfoID] = true
		}

		// 응답 데이터 생성
		songs := make([]SongSearchInfoV2Response, 0, len(songsWithArtist))
		for _, song := range songsWithArtist {
			songs = append(songs, SongSearchInfoV2Response{
				SongNumber:        song.SongNumber,
				SongName:          song.SongName,
				SingerName:        song.ArtistName,
				SongInfoId:        song.SongInfoID,
				Album:             song.Album.String,
				IsMr:              song.IsMR.Bool,
				IsLive:            song.IsLive.Bool,
				MelonLink:         CreateMelonLinkByMelonSongId(song.MelonSongID),
				IsKeep:            keepSongMap[song.SongInfoID],
				LyricsYoutubeLink: song.LyricsVideoLink.String,
				TJYoutubeLink:     song.TJYoutubeLink.String,
			})
		}
		response := SongSearchPageV2Response{
			Songs:    songs,
			NextPage: page + 1,
		}
		// 응답 반환
		pkg.BaseResponse(c, http.StatusOK, "ok", response)
	}
}

// SearchSongsBySongNameV2 godoc
// @Summary      노래 제목으로 노래 검색 API V2
// @Description  노래 제목으로 노래 검색 API V2, 노래 제목을 검색합니다. \n 검색 결과는 노래 제목, 아티스트 이름, 앨범명, 노래 번호 및 IsKeep여부를 반환합니다.
// @Tags         Search
// @Accept       json
// @Produce      json
// @Param        keyword query string true "검색 키워드"
// @Param        page query int false "현재 조회할 노래 목록의 쪽수. 입력하지 않는다면 기본값인 1쪽을 조회"
// @Param        size query int false "한번에 조회할 노래 개수. 입력하지 않는다면 기본값인 20개씩 조회"
// @Success      200 {object} pkg.BaseResponseStruct{data=SongSearchPageV2Response} "성공"
// @Failure      400 {object} pkg.BaseResponseStruct{data=nil} "실패"
// @Router       /v2/search/song-name [get]
// @Security BearerAuth
func SearchSongsBySongNameV2(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		memberId, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - memberId not found", nil)
			return
		}
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
			qm.OrderBy("CASE WHEN song_name LIKE ? THEN 1 WHEN song_name LIKE ? THEN 2 ELSE 3 END, song_info_id", searchKeyword, searchKeyword+"%"),
			qm.Limit(size),
			qm.Offset(offset),
		).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		keepList, err := mysql.KeepLists(
			qm.Where("member_id = ?", memberId),
		).One(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		keepSongs, err := mysql.KeepSongs(
			qm.Where("keep_list_id = ?", keepList.KeepListID),
		).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		keepSongMap := make(map[int64]bool)
		for _, keepSong := range keepSongs {
			keepSongMap[keepSong.SongInfoID] = true
		}

		// 응답 데이터 생성
		songs := make([]SongSearchInfoV2Response, 0, len(songsWithName))
		for _, song := range songsWithName {
			songs = append(songs, SongSearchInfoV2Response{
				SongNumber:        song.SongNumber,
				SongName:          song.SongName,
				SingerName:        song.ArtistName,
				SongInfoId:        song.SongInfoID,
				Album:             song.Album.String,
				IsMr:              song.IsMR.Bool,
				IsLive:            song.IsLive.Bool,
				MelonLink:         CreateMelonLinkByMelonSongId(song.MelonSongID),
				IsKeep:            keepSongMap[song.SongInfoID],
				LyricsYoutubeLink: song.LyricsVideoLink.String,
				TJYoutubeLink:     song.TJYoutubeLink.String,
			})
		}
		response := SongSearchPageV2Response{
			Songs:    songs,
			NextPage: page + 1,
		}
		// 응답 반환
		pkg.BaseResponse(c, http.StatusOK, "ok", response)
	}
}

// SearchSongsBySongNumberV2 godoc
// @Summary      노래 번호로 노래 검색 API V2
// @Description  노래 번호로 노래 검색 API V2, 노래 번호를 검색합니다. \n 검색 결과는 노래 제목, 아티스트 이름, 앨범명, 노래 번호및 IsKeep여부를 반환합니다.
// @Tags         Search
// @Accept       json
// @Produce      json
// @Param        keyword query string true "검색 키워드"
// @Param        page query int false "현재 조회할 노래 목록의 쪽수. 입력하지 않는다면 기본값인 1쪽을 조회. 현재는 노래 번호가 정확히 일치하는 1개만 반환하기 때문에 무의미"
// @Param        size query int false "한번에 조회할 노래 개수. 입력하지 않는다면 기본값인 20개씩 조회. 현재는 노래 번호가 정확히 일치하는 1개만 반환하기 때문에 무의미"
// @Success      200 {object} pkg.BaseResponseStruct{data=SongSearchPageV2Response} "성공"
// @Failure      400 {object} pkg.BaseResponseStruct{data=nil} "실패"
// @Router       /v2/search/song-number [get]
// @Security BearerAuth
func SearchSongsBySongNumberV2(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		memberId, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - memberId not found", nil)
			return
		}

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

		keepList, err := mysql.KeepLists(
			qm.Where("member_id = ?", memberId),
		).One(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		keepSongs, err := mysql.KeepSongs(
			qm.Where("keep_list_id = ?", keepList.KeepListID),
		).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		keepSongMap := make(map[int64]bool)
		for _, keepSong := range keepSongs {
			keepSongMap[keepSong.SongInfoID] = true
		}

		// 응답 데이터 생성
		songs := make([]SongSearchInfoV2Response, 0, len(songsWithNumber))
		for _, song := range songsWithNumber {
			songs = append(songs, SongSearchInfoV2Response{
				SongNumber:        song.SongNumber,
				SongName:          song.SongName,
				SingerName:        song.ArtistName,
				SongInfoId:        song.SongInfoID,
				Album:             song.Album.String,
				IsMr:              song.IsMR.Bool,
				IsLive:            song.IsLive.Bool,
				MelonLink:         CreateMelonLinkByMelonSongId(song.MelonSongID),
				IsKeep:            keepSongMap[song.SongInfoID],
				LyricsYoutubeLink: song.LyricsVideoLink.String,
				TJYoutubeLink:     song.TJYoutubeLink.String,
			})
		}
		response := SongSearchPageV2Response{
			Songs:    songs,
			NextPage: page + 1,
		}
		// 응답 반환
		pkg.BaseResponse(c, http.StatusOK, "ok", response)
	}
}
