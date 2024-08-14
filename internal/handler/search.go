package handler

import (
	"SingSong-Server/internal/db/mysql"
	"SingSong-Server/internal/pkg"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"net/http"
)

// getUnicodeRange returns the Unicode range for a given Korean consonant.
func getUnicodeRange(consonant string) (string, string, bool) {
	ranges := map[string][2]string{
		"ㄱ": {"가", "깋"},
		"ㄴ": {"나", "닣"},
		"ㄷ": {"다", "딯"},
		"ㄹ": {"라", "맇"},
		"ㅁ": {"마", "밓"},
		"ㅂ": {"바", "빟"},
		"ㅅ": {"사", "싷"},
		"ㅇ": {"아", "잏"},
		"ㅈ": {"자", "짛"},
		"ㅊ": {"차", "칳"},
		"ㅋ": {"카", "킿"},
		"ㅌ": {"타", "팋"},
		"ㅍ": {"파", "핗"},
		"ㅎ": {"하", "힣"},
	}

	val, exists := ranges[consonant]
	if !exists {
		return "", "", false
	}
	return val[0], val[1], true
}

// SearchSongs godoc
// @Summary      노래 검색 API
// @Description  노래 검색 API로, 노래 제목 또는 아티스트 이름을 검색합니다.
// @Tags         Search
// @Accept       json
// @Produce      json
// @Param        searchKeyword path string true "검색 키워드"
// @Success      200 {object} pkg.BaseResponseStruct{data=map[string][]songInfoResponse} "성공"
// @Failure      400 {object} pkg.BaseResponseStruct{data=map[string][]songInfoResponse} "실패 - 빈 리스트 반환"
// @Router       /search/{searchKeyword} [get]
func SearchSongs(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 검색어를 URL 파라미터에서 가져오기
		searchKeyword := c.Param("searchKeyword")

		// 노래 이름으로 검색
		songsWithName, err := mysql.SongInfos(
			qm.Where("song_name LIKE ?", "%"+searchKeyword+"%"),
			qm.Limit(10),
		).All(c, db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// 아티스트 이름으로 검색
		songsWithArtist, err := mysql.SongInfos(
			qm.Where("artist_name LIKE ?", "%"+searchKeyword+"%"),
			qm.Limit(10),
		).All(c, db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// 검색 결과를 담을 맵 생성
		response := make(map[string][]songInfoResponse)

		// songName 해당하는 검색 결과를 response 추가
		response["songName"] = make([]songInfoResponse, len(songsWithName))
		for i, song := range songsWithName {
			response["songName"][i] = songInfoResponse{
				SongNumber:  song.SongNumber,
				SongName:    song.SongName,
				SingerName:  song.ArtistName,
				Tags:        parseTags(song.Tags.String),
				SongInfoId:  song.SongInfoID,
				Album:       song.Album.String,
				Octave:      song.Octave.String,
				Description: "", // todo: Add description logic
			}
		}

		// artistName 해당하는 검색 결과를 response 추가
		response["artistName"] = make([]songInfoResponse, len(songsWithArtist))
		for i, song := range songsWithArtist {
			response["artistName"][i] = songInfoResponse{
				SongNumber:  song.SongNumber,
				SongName:    song.SongName,
				SingerName:  song.ArtistName,
				Tags:        parseTags(song.Tags.String),
				SongInfoId:  song.SongInfoID,
				Album:       song.Album.String,
				Octave:      song.Octave.String,
				Description: "", // todo: Add description logic
			}
		}

		// 응답 반환
		pkg.BaseResponse(c, http.StatusOK, "ok", response)
	}
}
