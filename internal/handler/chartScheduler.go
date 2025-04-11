package handler

//
//import (
//	"context"
//	"database/sql"
//	"encoding/json"
//	"fmt"
//	"github.com/pkg/errors"
//	"github.com/redis/go-redis/v9"
//	"github.com/volatiletech/null/v8"
//	"log"
//	"sort"
//	"strings"
//	"time"
//)
//
//type SongStruct struct {
//	SongInfoID      int         `json:"song_info_id"`
//	TotalScore      float64     `json:"total_score"`
//	SongName        string      `json:"song_name"`
//	ArtistName      string      `json:"artist_name"`
//	SongNumber      int         `json:"song_number"`
//	IsMR            bool        `json:"is_mr"`
//	IsLive          bool        `json:"is_live"`
//	Album           null.String `json:"album,omitempty"`
//	Gender          null.String `json:"gender"`
//	AgeGroup        string      `json:"age_group"`
//	MelonSongID     null.String `json:"melon_song_id,omitempty"`
//	LyricsVideoLink null.String `json:"lyrics_video_link,omitempty"`
//	TJYoutubeLink   null.String `json:"tj_youtube_link,omitempty"`
//	Ranking         int         `json:"ranking"`
//	RankingChange   int         `json:"ranking_change"`
//	New             bool        `json:"new"`
//}
//
//func ScheduleNextChart(db *sql.DB, rdb *redis.Client) {
//	log.Println("chart cronjob start!")
//	history, err := queryHistory(db)
//	if err != nil {
//		log.Printf(err.Error())
//		return
//	}
//
//	// 상호작용 데이터를 성별/연령대/노래에 따라 분류
//	maleData := map[string][]SongStruct{"ALL": {}, "10": {}, "20": {}, "30": {}, "40+": {}}
//	femaleData := map[string][]SongStruct{"ALL": {}, "10": {}, "20": {}, "30": {}, "40+": {}}
//	mixedData := map[string][]SongStruct{"ALL": {}, "10": {}, "20": {}, "30": {}, "40+": {}}
//
//	for _, row := range history {
//		mixedData["ALL"] = updateData(mixedData["ALL"], row, "MIXED", "ALL")
//		if row.AgeGroup != "ALL" {
//			mixedData[row.AgeGroup] = updateData(mixedData[row.AgeGroup], row, "MIXED", row.AgeGroup)
//		}
//		if row.Gender.String == "MALE" {
//			maleData["ALL"] = updateData(maleData["ALL"], row, "MALE", "ALL")
//			if row.AgeGroup != "ALL" {
//				maleData[row.AgeGroup] = updateData(maleData[row.AgeGroup], row, "MALE", row.AgeGroup)
//			}
//		} else if row.Gender.String == "FEMALE" {
//			femaleData["ALL"] = updateData(femaleData["ALL"], row, "FEMALE", "ALL")
//			if row.AgeGroup != "ALL" {
//				femaleData[row.AgeGroup] = updateData(femaleData[row.AgeGroup], row, "FEMALE", row.AgeGroup)
//			}
//		}
//	}
//
//	ctx := context.Background()
//	previousData, err := fetchPreviousData(ctx, rdb)
//	if err != nil {
//		log.Printf(err.Error())
//		return
//	}
//
//	// 상위 20개의 노래를 뽑아서 이전 데이터와 비교해서 랭킹 변화 정보를 추가
//	for ageGroup := range maleData {
//		maleData[ageGroup] = addRankingInfo(getTop20ByScore(maleData[ageGroup]), previousData["MALE"][ageGroup])
//		femaleData[ageGroup] = addRankingInfo(getTop20ByScore(femaleData[ageGroup]), previousData["FEMALE"][ageGroup])
//		mixedData[ageGroup] = addRankingInfo(getTop20ByScore(mixedData[ageGroup]), previousData["MIXED"][ageGroup])
//	}
//
//	now := time.Now()
//	baseKey := now.Add(time.Hour).Format("2006-01-02-15") + "-Hot_Trend"
//	if err := saveToRedis(ctx, rdb, baseKey, maleData, "MALE"); err != nil {
//		log.Printf("Failed to save male data to Redis: %v", err)
//	}
//	if err := saveToRedis(ctx, rdb, baseKey, femaleData, "FEMALE"); err != nil {
//		log.Printf("Failed to save female data to Redis: %v", err)
//	}
//	if err := saveToRedis(ctx, rdb, baseKey, mixedData, "MIXED"); err != nil {
//		log.Printf("Failed to save mixed data to Redis: %v", err)
//	}
//	log.Println("Hot trending data updated successfully.")
//}
//
//func queryHistory(db *sql.DB) ([]SongStruct, error) {
//	query := fmt.Sprintf(`
//		WITH scored_songs AS (
//			SELECT
//				ma.song_info_id,
//				SUM(ma.action_score) AS total_score,
//				ma.gender,
//				CASE
//					WHEN ma.birthyear IS NULL OR ma.birthyear = 0 OR ma.birthyear > YEAR(CURDATE()) THEN 'ALL'
//					WHEN YEAR(CURDATE()) - ma.birthyear + 1 BETWEEN 10 AND 19 THEN '10'
//					WHEN YEAR(CURDATE()) - ma.birthyear + 1 BETWEEN 20 AND 29 THEN '20'
//					WHEN YEAR(CURDATE()) - ma.birthyear + 1 BETWEEN 30 AND 39 THEN '30'
//					WHEN YEAR(CURDATE()) - ma.birthyear + 1 > 39 THEN '40+'
//					ELSE 'ALL'
//				END AS age_group
//			FROM member_action AS ma
//			WHERE ma.created_at > DATE_SUB(NOW(), INTERVAL 8 WEEK)
//			GROUP BY ma.song_info_id, ma.gender, age_group
//		)
//		SELECT
//			ss.song_info_id,
//			ss.total_score,
//			s.song_name,
//			s.artist_name,
//			s.song_number,
//			s.is_mr,
//			s.is_live,
//			s.album,
//			ss.gender,
//			ss.age_group,
//			s.melon_song_id,
//			s.lyrics_video_link,
//			s.tj_youtube_link
//		FROM scored_songs ss
//		JOIN song_info s ON ss.song_info_id = s.song_info_id
//	`)
//
//	rows, err := db.Query(query)
//	if err != nil {
//		return nil, errors.Wrap(err, "인기차트 생성을 위한 쿼리 조회 실패")
//	}
//	defer rows.Close()
//
//	results := make([]SongStruct, 0)
//	count := 0
//	for rows.Next() {
//		count++
//		var song SongStruct
//		err := rows.Scan(
//			&song.SongInfoID,
//			&song.TotalScore,
//			&song.SongName,
//			&song.ArtistName,
//			&song.SongNumber,
//			&song.IsMR,
//			&song.IsLive,
//			&song.Album,
//			&song.Gender,
//			&song.AgeGroup,
//			&song.MelonSongID,
//			&song.LyricsVideoLink,
//			&song.TJYoutubeLink,
//		)
//
//		if err != nil {
//			return nil, errors.Wrap(err, "인기차트 생성을 위한 쿼리 조회 실패")
//		}
//		results = append(results, song)
//	}
//	return results, nil
//}
//
//func updateData(dataList []SongStruct, row SongStruct, gender, ageGroup string) []SongStruct {
//	for i := range dataList {
//		if (dataList)[i].SongInfoID == row.SongInfoID {
//			(dataList)[i].TotalScore += row.TotalScore
//			return dataList
//		}
//	}
//	row.Gender = null.StringFrom(gender)
//	row.AgeGroup = ageGroup
//	dataList = append(dataList, row)
//	return dataList
//}
//
//func getTop20ByScore(dataList []SongStruct) []SongStruct {
//	sort.Slice(dataList, func(i, j int) bool {
//		return dataList[i].TotalScore > dataList[j].TotalScore
//	})
//	if len(dataList) > 20 {
//		return dataList[:20]
//	}
//	return dataList
//}
//
//func addRankingInfo(currentData, previousData []SongStruct) []SongStruct {
//	previousRanking := make(map[int]int)
//	for _, item := range previousData {
//		previousRanking[item.SongInfoID] = item.Ranking
//	}
//
//	for idx, item := range currentData {
//		item.Ranking = idx + 1
//		if prevRank, exists := previousRanking[item.SongInfoID]; exists {
//			item.RankingChange = prevRank - item.Ranking
//			item.New = false
//		} else {
//			item.RankingChange = 0
//			item.New = true
//		}
//		currentData[idx] = item
//	}
//	return currentData
//}
//
//func saveToRedis(ctx context.Context, redisClient *redis.Client, baseKey string, data map[string][]SongStruct, gender string) error {
//	for ageGroup, items := range data {
//		key := fmt.Sprintf("%s_%s_%s", baseKey, gender, ageGroup)
//		jsonData, err := json.Marshal(items)
//		if err != nil {
//			return errors.Wrap(err, "최초 에러 발생 지점")
//		}
//		if err := redisClient.Set(ctx, key, jsonData, 4800*time.Second).Err(); err != nil {
//			return errors.Wrap(err, "최초 에러 발생 지점")
//		}
//	}
//	return nil
//}
//
//func fetchPreviousData(ctx context.Context, redisClient *redis.Client) (map[string]map[string][]SongStruct, error) {
//	// 현재 시간과 키 포맷 설정
//	now := time.Now()
//	formattedKeyBase := now.Format("2006-01-02-15-Hot_Trend")
//
//	// 결과 데이터 초기화
//	previousData := map[string]map[string][]SongStruct{
//		"male":   {},
//		"female": {},
//		"mixed":  {},
//	}
//
//	// 성별과 연령대별 데이터를 가져오기
//	genderKeys := []string{"MALE", "FEMALE", "MIXED"}
//	ageGroups := []string{"ALL", "10", "20", "30", "40+"}
//
//	for _, genderKey := range genderKeys {
//		gender := strings.ToLower(genderKey)
//		previousData[gender] = make(map[string][]SongStruct)
//
//		for _, ageGroup := range ageGroups {
//			// Redis 키 생성
//			key := fmt.Sprintf("%s_%s_%s", formattedKeyBase, genderKey, ageGroup)
//
//			// Redis에서 데이터 가져오기
//			data, err := redisClient.Get(ctx, key).Result()
//			if err == redis.Nil {
//				// 키가 없을 경우 빈 리스트로 초기화
//				previousData[gender][ageGroup] = []SongStruct{}
//				continue
//			} else if err != nil {
//				log.Printf("Error fetching key %s: %v", key, err)
//				return nil, errors.Wrap(err, "최초 에러 발생 지점")
//			}
//
//			// JSON 파싱
//			var parsedData []SongStruct
//			if err := json.Unmarshal([]byte(data), &parsedData); err != nil {
//				log.Printf("Invalid JSON for key %s. Setting as empty list. Error: %v", key, err)
//				previousData[gender][ageGroup] = []SongStruct{}
//				continue
//			}
//
//			// 유효한 데이터 저장
//			previousData[gender][ageGroup] = parsedData
//		}
//	}
//
//	return previousData, nil
//}
//
//// 서버 실행 시 해당 시각의 인기 차트가 없을 경우 생성하도록 하는 코드
//func InitializeChart(db *sql.DB, rdb *redis.Client) {
//	log.Printf("chart initialization start!")
//	ctx := context.Background()
//
//	now := time.Now()
//	currentInitKey := now.Format("2006-01-02-15") + "-Hot_Trend_MIXED_ALL"
//	nextInitKey := now.Add(time.Hour).Format("2006-01-02-15") + "-Hot_Trend_MIXED_ALL"
//
//	// Redis에서 현재와 다음 키 존재 여부 확인
//	currentExists, err := rdb.Exists(ctx, currentInitKey).Result()
//	if err != nil {
//		log.Printf("Error checking current key in Redis: %v", err)
//	}
//	nextExists, err := rdb.Exists(ctx, nextInitKey).Result()
//	if err != nil {
//		log.Printf("Error checking next key in Redis: %v", err)
//	}
//
//	// 없으면 추가한다
//	if currentExists == 0 {
//		log.Printf("Initializing current chart: %s", currentInitKey)
//		scheduleCurrentChart(db, rdb)
//	}
//
//	if nextExists == 0 {
//		log.Printf("Initializing next chart: %s", nextInitKey)
//		ScheduleNextChart(db, rdb)
//	}
//
//	log.Printf("chart initialization end!")
//}
//
//func scheduleCurrentChart(db *sql.DB, rdb *redis.Client) {
//	history, err := queryHistory(db)
//	if err != nil {
//		log.Printf(err.Error())
//		return
//	}
//
//	// 상호작용 데이터를 성별/연령대/노래에 따라 분류
//	maleData := map[string][]SongStruct{"ALL": {}, "10": {}, "20": {}, "30": {}, "40+": {}}
//	femaleData := map[string][]SongStruct{"ALL": {}, "10": {}, "20": {}, "30": {}, "40+": {}}
//	mixedData := map[string][]SongStruct{"ALL": {}, "10": {}, "20": {}, "30": {}, "40+": {}}
//
//	for _, row := range history {
//		mixedData["ALL"] = updateData(mixedData["ALL"], row, "MIXED", "ALL")
//		if row.AgeGroup != "ALL" {
//			mixedData[row.AgeGroup] = updateData(mixedData[row.AgeGroup], row, "MIXED", row.AgeGroup)
//		}
//		if row.Gender.String == "MALE" {
//			maleData["ALL"] = updateData(maleData["ALL"], row, "MALE", "ALL")
//			if row.AgeGroup != "ALL" {
//				maleData[row.AgeGroup] = updateData(maleData[row.AgeGroup], row, "MALE", row.AgeGroup)
//			}
//		} else if row.Gender.String == "FEMALE" {
//			femaleData["ALL"] = updateData(femaleData["ALL"], row, "FEMALE", "ALL")
//			if row.AgeGroup != "ALL" {
//				femaleData[row.AgeGroup] = updateData(femaleData[row.AgeGroup], row, "FEMALE", row.AgeGroup)
//			}
//		}
//	}
//
//	ctx := context.Background()
//	previousData := map[string]map[string][]SongStruct{
//		"male":   {},
//		"female": {},
//		"mixed":  {},
//	}
//
//	// 상위 20개의 노래를 뽑아서 이전 데이터와 비교해서 랭킹 변화 정보를 추가
//	for ageGroup := range maleData {
//		maleData[ageGroup] = addRankingInfo(getTop20ByScore(maleData[ageGroup]), previousData["MALE"][ageGroup])
//		femaleData[ageGroup] = addRankingInfo(getTop20ByScore(femaleData[ageGroup]), previousData["FEMALE"][ageGroup])
//		mixedData[ageGroup] = addRankingInfo(getTop20ByScore(mixedData[ageGroup]), previousData["MIXED"][ageGroup])
//	}
//
//	now := time.Now()
//	baseKey := now.Format("2006-01-02-15") + "-Hot_Trend"
//	if err := saveToRedis(ctx, rdb, baseKey, maleData, "MALE"); err != nil {
//		log.Printf("Failed to save male data to Redis: %v", err)
//	}
//	if err := saveToRedis(ctx, rdb, baseKey, femaleData, "FEMALE"); err != nil {
//		log.Printf("Failed to save female data to Redis: %v", err)
//	}
//	if err := saveToRedis(ctx, rdb, baseKey, mixedData, "MIXED"); err != nil {
//		log.Printf("Failed to save mixed data to Redis: %v", err)
//	}
//}
