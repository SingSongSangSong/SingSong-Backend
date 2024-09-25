package handler

import (
	"SingSong-Server/conf"
	"SingSong-Server/internal/db/mysql"
	"SingSong-Server/internal/pkg"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"log"
	"net/http"
	"strconv"
)

// RelatedSongV2 godoc
// @Summary      연관된 노래들을 조회합니다
// @Description  연관된 노래들과 다음 페이지 번호를 함께 조회합니다. 노래 상세 화면에 첫 진입했을 경우 page 번호는 1입니다. 무한스크롤을 진행한다면 응답에 포함되어 오는 nextPage를 다음번에 포함하여 보내면 됩니다. nextPage는 1씩 증가합니다. 더이상 노래가 없을 경우, 응답에는 빈 배열과 함께 nextPage는 1로 반환됩니다.
// @Tags         Songs
// @Accept       json
// @Produce      json
// @Param        songId path string true "songId"
// @Param        page query int false "현재 조회할 노래 목록의 쪽수. 입력하지 않는다면 기본값인 1쪽을 조회"
// @Param        size query int false "한번에 조회할 노래 개수. 입력하지 않는다면 기본값인 20개씩 조회"
// @Success      200 {object} pkg.BaseResponseStruct{data=relatedSongResponse} "성공"
// @Router       /v2/songs/{songId}/related [get]
// @Security BearerAuth
func RelatedSongV2(db *sql.DB, milvusClient *client.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. songInfoId 및 memberId 검증
		songInfoId := c.Param("songId")
		if songInfoId == "" {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - cannot find songId in path variable", nil)
			return
		}

		value, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - memberId not found", nil)
			return
		}

		memberId, ok := value.(int64)
		if !ok {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - memberId not type int64", nil)
			return
		}

		// 2. 페이지 및 크기 값 설정
		sizeStr := c.DefaultQuery("size", defaultSize)
		sizeInt, err := strconv.Atoi(sizeStr)
		if err != nil || sizeInt < 0 {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid size parameter", nil)
			return
		}

		pageStr := c.DefaultQuery("page", defaultPage)
		pageInt, err := strconv.Atoi(pageStr)
		if err != nil || pageInt < 0 {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid page parameter", nil)
			return
		}

		// 3. 벡터 크기 및 마지막 페이지 여부 설정
		vectorSize := sizeInt * pageInt
		isLastPage := false
		if vectorSize > maximumSongSize && vectorSize-sizeInt < maximumSongSize {
			vectorSize = maximumSongSize
			isLastPage = true
		} else if vectorSize > maximumSongSize && vectorSize-sizeInt >= maximumSongSize {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - related song data limit is 100", nil)
			return
		}

		// 다음과 같이 Get함수를 사용하여 song_info_id를 가져오는 방법은 사용하지 않습니다.
		// 그 이유는 TopK, Offset, Limit등 다양한 설정이 불가하기 때문입니다.
		//// 1. song_info_id 값을 담는 ColumnString 생성
		//songInfoIdForGet := []string{songInfoId} // 실제 song_info_id 리스트로 변경
		//columnString := entity.NewColumnString("song_info_id", songInfoIdForGet)
		//
		//// Milvus에서 Get 함수 호출
		//resultSet, err := (*milvusClient).Get(c, conf.VectorDBConfigInstance.COLLECTION_NAME, columnString)
		//if err != nil {
		//	log.Fatalf("Milvus에서 데이터를 가져오는 데 실패했습니다: %v", err)
		//}

		// 4. 벡터 디비에서 조회
		songVector, err := (*milvusClient).Query(
			c,
			conf.VectorDBConfigInstance.COLLECTION_NAME,
			[]string{},
			"song_info_id == "+songInfoId,
			[]string{"vector"},
		)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		floatValues := songVector.GetColumn("vector").FieldData().GetVectors().GetFloatVector().GetData()
		if len(floatValues) == 0 {
			pkg.BaseResponse(c, http.StatusOK, "ok", relatedSongResponse{[]relatedSong{}, 1})
			return
		}
		vectorData := entity.FloatVector(floatValues)

		// 6. 벡터 기반으로 연관 노래 검색
		sp, _ := entity.NewIndexFlatSearchParam()
		sr, err := (*milvusClient).Search(
			c,
			conf.VectorDBConfigInstance.COLLECTION_NAME,
			[]string{},
			"song_info_id != "+songInfoId+" && MR == false",
			[]string{"song_name", "artist_name", "album", "song_number", "MR", "song_info_id"},
			[]entity.Vector{vectorData},
			"vector",
			entity.COSINE,
			sizeInt,
			sp,
			client.WithOffset(int64(pageInt)),
		)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}
		if len(sr) == 0 {
			pkg.BaseResponse(c, http.StatusOK, "ok", relatedSongResponse{[]relatedSong{}, 1})
			return
		}

		// 7. 사용자가 보관한 노래 조회 (isKeep)
		all, err := mysql.KeepLists(qm.Where("member_id = ?", memberId)).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}
		keepIds := make([]interface{}, len(all))
		for i, keep := range all {
			keepIds[i] = keep.KeepListID
		}

		keepSongs, err := mysql.KeepSongs(qm.WhereIn("keep_list_id in ?", keepIds...), qm.And("deleted_at IS NULL")).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		isKeepMap := make(map[int64]bool)
		for _, keepSong := range keepSongs {
			isKeepMap[keepSong.SongInfoID] = true
		}

		// 8. songInfoId로 노래 정보 가져오기
		var songInfoIds []interface{}
		for _, song := range sr {
			// song.Fields.GetColumn("song_info_id")가 nil인지 확인
			column := song.Fields.GetColumn("song_info_id")
			if column == nil {
				log.Printf("song_info_id 컬럼을 찾을 수 없습니다. song: %v", song.Fields)
				continue // 에러 발생 시 다음 song으로 건너뜁니다.
			}

			// FieldData가 nil인지 확인
			fieldData := column.FieldData()
			if fieldData == nil {
				log.Printf("song_info_id의 FieldData를 찾을 수 없습니다. song: %v", song)
				continue // 에러 발생 시 다음 song으로 건너뜁니다.
			}

			// 실제 song_info_id 값을 추출
			longData := fieldData.GetScalars().GetLongData().GetData()
			if len(longData) == 0 {
				log.Printf("song_info_id의 데이터를 찾을 수 없습니다. song: %v", fieldData)
				continue // 에러 발생 시 다음 song으로 건
			}
			for _, val := range longData {
				songInfoIds = append(songInfoIds, val)
			}
		}

		// songInfoIds에 값이 있는지 확인
		if len(songInfoIds) == 0 {
			pkg.BaseResponse(c, http.StatusNotFound, "관련된 노래 정보를 찾을 수 없습니다.", nil)
			return
		}

		slice, err := mysql.SongInfos(qm.WhereIn("song_info_id IN ?", songInfoIds...)).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		songInfoMap := make(map[int64]*mysql.SongInfo)
		for _, song := range slice {
			songInfoMap[song.SongInfoID] = song
		}

		// 9. 관련 노래 목록 생성 및 응답 생성
		relatedSongs := make([]relatedSong, 0, sizeInt)
		for _, song := range songInfoIds {
			found := songInfoMap[song.(int64)]
			relatedSongs = append(relatedSongs, relatedSong{
				SongInfoId: found.SongInfoID,
				SongName:   found.SongName,
				SingerName: found.ArtistName,
				Album:      found.Album.String,
				IsKeep:     isKeepMap[found.SongInfoID],
				SongNumber: found.SongNumber,
				IsMr:       found.IsMR.Bool,
				IsLive:     found.IsLive.Bool,
				MelonLink:  CreateMelonLinkByMelonSongId(found.MelonSongID),
			})
		}
		nextPage := pageInt + 1
		if isLastPage {
			nextPage = 1
		}
		pkg.BaseResponse(c, http.StatusOK, "ok", relatedSongResponse{relatedSongs, nextPage})
	}
}
