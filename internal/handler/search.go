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

func SearchSongs(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		memberId := c.Param("memberId")
		if memberId == "" {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - cannot find memberId in path variable", nil)
			return
		}

		searchKeyword := c.Param("searchKeyword")
		if searchKeyword == "" {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - cannot find searchKeyword in path variable", nil)
			return
		}

		var one []*mysql.SongInfo
		var err error

		// 초성인지 확인
		if startChar, endChar, ok := getUnicodeRange(searchKeyword); ok {
			// 초성 검색 (song_name 및 artist_name)
			one, err = mysql.SongInfos(
				qm.Where("(song_name >= ? AND song_name <= ?) OR (artist_name >= ? AND artist_name <= ?)", startChar, endChar, startChar, endChar),
			).All(c, db)
		} else if _, err := strconv.Atoi(searchKeyword); err == nil {
			// 숫자 검색 (song_number)
			one, err = mysql.SongInfos(
				qm.Where("song_number = ?", searchKeyword),
			).All(c, db)
		} else {
			// 일반 검색 (song_name 및 artist_name)
			one, err = mysql.SongInfos(
				qm.Where("song_name LIKE ? OR artist_name LIKE ?", "%"+searchKeyword+"%", "%"+searchKeyword+"%"),
			).All(c, db)
		}

		if err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - no song", nil)
			return
		}

		// 유저의 keep 여부 조회
		all, err := mysql.KeepLists(
			qm.Where("member_id = ?", memberId),
		).All(c, db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// 결과
		response := make([]songInfoResponse, len(one))
		for i, song := range one {
			keepListIds := make([]interface{}, len(all))
			for j, keep := range all {
				keepListIds[j] = keep.KeepListID
			}
			isKeep, err := mysql.KeepSongs(
				qm.WhereIn("keep_list_id in ?", keepListIds...),
				qm.And("song_info_id = ?", song.SongInfoID),
				qm.And("deleted_at IS NULL"),
			).Exists(c, db)
			if err != nil {
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}

			response[i] = songInfoResponse{
				SongNumber:  song.SongNumber,
				SongName:    song.SongName,
				SingerName:  song.ArtistName,
				Tags:        parseTags(song.Tags.String),
				SongInfoId:  song.SongInfoID,
				Album:       song.Album.String,
				Octave:      song.Octave.String,
				Description: "", //todo:
				IsKeep:      isKeep,
			}
		}

		pkg.BaseResponse(c, http.StatusOK, "ok", response)
	}
}
