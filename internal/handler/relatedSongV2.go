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
		if vectorSize > maximumSongSize {
			if vectorSize-sizeInt < maximumSongSize {
				vectorSize = maximumSongSize
				isLastPage = true
			} else {
				pkg.BaseResponse(c, http.StatusBadRequest, "error - related song data limit is 100", nil)
				return
			}
		}

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

		// 5. 벡터 데이터 추출
		var vectorData entity.FloatVector
		for _, column := range songVector {
			if column.Name() == "vector" {
				fieldData := column.FieldData()

				if fieldData == nil || fieldData.GetVectors() == nil {
					log.Fatalf("벡터 데이터를 찾을 수 없습니다.")
				}

				// 벡터 데이터를 FloatVector로 변환
				floatValues := fieldData.GetVectors().GetFloatVector().GetData()
				vectorData = entity.FloatVector(floatValues)
				break
			}
		}

		// 6. 벡터 기반으로 연관 노래 검색
		sp, _ := entity.NewIndexFlatSearchParam()
		sr, err := (*milvusClient).Search(
			c,
			conf.VectorDBConfigInstance.COLLECTION_NAME,
			[]string{},
			"song_info_id != "+songInfoId,
			[]string{"song_name", "artist_name", "album", "song_number", "MR"},
			[]entity.Vector{vectorData},
			"vector",
			entity.COSINE,
			20,
			sp,
			client.WithOffset(int64(pageInt)),
			client.WithLimit(int64(vectorSize)),
		)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
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
			songInfoIds = append(songInfoIds, song.Fields.GetColumn("song_info_id").FieldData())
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
			})
		}
	}
}
