package handler

import (
	"SingSong-Server/internal/db/mysql"
	"SingSong-Server/internal/pkg"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/pinecone-io/go-pinecone/pinecone"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"google.golang.org/protobuf/types/known/structpb"
	"net/http"
	"strconv"
)

type relatedSong struct {
	SongNumber        int    `json:"songNumber"`
	SongName          string `json:"songName"`
	SingerName        string `json:"singerName"`
	Album             string `json:"album"`
	IsKeep            bool   `json:"isKeep"`
	SongInfoId        int64  `json:"songId"`
	IsMr              bool   `json:"isMr"`
	IsLive            bool   `json:"isLive"`
	MelonLink         string `json:"melonLink"`
	LyricsYoutubeLink string `json:"lyricsYoutubeLink"`
	TJYoutubeLink     string `json:"tjYoutubeLink"`
}

type relatedSongResponse struct {
	Songs    []relatedSong `json:"songs"`
	NextPage int           `json:"nextPage"`
}

var (
	defaultSize     = "20"
	defaultPage     = "1"
	maximumSongSize = 100
)

// RelatedSong godoc
// @Summary      연관된 노래들을 조회합니다
// @Description  연관된 노래들과 다음 페이지 번호를 함께 조회합니다. 노래 상세 화면에 첫 진입했을 경우 page 번호는 1입니다. 무한스크롤을 진행한다면 응답에 포함되어 오는 nextPage를 다음번에 포함하여 보내면 됩니다. nextPage는 1씩 증가합니다. 더이상 노래가 없을 경우, 응답에는 빈 배열과 함께 nextPage는 1로 반환됩니다.
// @Tags         Songs
// @Accept       json
// @Produce      json
// @Param        songId path string true "songId"
// @Param        page query int false "현재 조회할 노래 목록의 쪽수. 입력하지 않는다면 기본값인 1쪽을 조회"
// @Param        size query int false "한번에 조회할 노래 개수. 입력하지 않는다면 기본값인 20개씩 조회"
// @Success      200 {object} pkg.BaseResponseStruct{data=relatedSongResponse} "성공"
// @Router       /v1/songs/{songId}/related [get]
// @Security BearerAuth
func RelatedSong(db *sql.DB, idxConnection *pinecone.IndexConnection) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		sizeStr := c.DefaultQuery("size", defaultSize)
		sizeInt, err := strconv.Atoi(sizeStr)
		if err != nil || sizeInt < 0 {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid size parameter", nil)
			return
		}

		pageStr := c.DefaultQuery("page", defaultPage)
		pageInt, err := strconv.Atoi(pageStr)
		if err != nil || pageInt < 0 {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid size parameter", nil)
			return
		}

		vectorSize := sizeInt * pageInt
		isLastPage := false
		if vectorSize > maximumSongSize && vectorSize-sizeInt < maximumSongSize {
			vectorSize = maximumSongSize
			isLastPage = true
		} else if vectorSize > maximumSongSize && vectorSize-sizeInt >= maximumSongSize {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - related song data limit is 100", nil)
			return
		}

		//songInfoId로 벡터 디비에서 조회
		songInfoIdInt, err := strconv.Atoi(songInfoId)
		if err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid songInfoId", nil)
			return
		}
		filterStruct := &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"song_id": structpb.NewStructValue(&structpb.Struct{
					Fields: map[string]*structpb.Value{
						"$ne": structpb.NewNumberValue(float64(songInfoIdInt)),
					},
				}),
				"MR": structpb.NewBoolValue(false),
			},
		}

		res, err := idxConnection.QueryByVectorId(c, &pinecone.QueryByVectorIdRequest{
			VectorId:        songInfoId,
			TopK:            uint32(vectorSize),
			Filter:          filterStruct,
			IncludeValues:   false,
			IncludeMetadata: false,
		})

		if len(res.Matches) == 0 {
			pkg.BaseResponse(c, http.StatusOK, "ok", relatedSongResponse{[]relatedSong{}, 1})
			return
		}
		res.Matches = res.Matches[sizeInt*(pageInt-1):]

		relatedSongs := make([]relatedSong, 0, sizeInt)
		if len(res.Matches) <= 0 {
			pkg.BaseResponse(c, http.StatusOK, "ok", relatedSongResponse{relatedSongs, 1})
			return
		}

		for _, each := range res.Matches {
			v := each.Vector
			atoi, err := strconv.Atoi(v.Id)
			if err != nil {
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}
			relatedSongs = append(relatedSongs, relatedSong{
				SongInfoId: int64(atoi),
			})
		}

		//isKeep
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

		//songInfoId로 노래 정보 가져오기
		var songInfoIds []interface{}
		for _, song := range relatedSongs {
			songInfoIds = append(songInfoIds, song.SongInfoId)
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

		// response에 isKeep과 songTempId 추가
		for i, song := range relatedSongs {
			found := songInfoMap[song.SongInfoId]
			relatedSongs[i].SongName = found.SongName
			relatedSongs[i].SingerName = found.ArtistName
			relatedSongs[i].Album = found.Album.String
			relatedSongs[i].IsKeep = isKeepMap[song.SongInfoId]
			relatedSongs[i].SongNumber = found.SongNumber
			relatedSongs[i].IsMr = found.IsMR.Bool
			relatedSongs[i].IsLive = found.IsLive.Bool
			relatedSongs[i].MelonLink = CreateMelonLinkByMelonSongId(found.MelonSongID)
			relatedSongs[i].LyricsYoutubeLink = found.LyricsVideoLink.String
			relatedSongs[i].TJYoutubeLink = found.TJYoutubeLink.String
		}

		nextPage := pageInt + 1
		if isLastPage {
			nextPage = 1
		}
		pkg.BaseResponse(c, http.StatusOK, "ok", relatedSongResponse{relatedSongs, nextPage})
	}
}
