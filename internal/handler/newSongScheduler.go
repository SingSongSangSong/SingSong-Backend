package handler

import (
	"SingSong-Server/internal/db/mysql"
	"context"
	"database/sql"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/agnivade/levenshtein"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type NewSongStruct struct {
	SongNumber  int
	SongName    string
	ArtistName  string
	IsMr        bool
	IsLive      bool
	MelonSongId string
	Album       string
	Genre       string
	ReleaseYear int
}

type SearchResult struct {
	SongName    string
	ArtistName  string
	MelonSongId string
}

func ScheduleNewSongs(db *sql.DB) {
	ctx := context.Background()

	// 신곡 가져오기
	newSongs, err := fetchNewSongs(db)
	if err != nil {
		log.Printf(err.Error())
		return
	}

	if newSongs == nil || len(newSongs) == 0 {
		return
	}

	// isMr과 isLive 가져와 저장
	updatedSongs, err := fetchMrAndLive(newSongs)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	err = saveToDB(ctx, db, updatedSongs)
	if err != nil {
		log.Printf("Error saving songs to DB: %v", err)
		return
	}

	// melon song id 가져오기
	updatedSongs, err = fetchMelonSongId(updatedSongs)
	if err != nil {
		log.Printf(err.Error())
		return
	}

	// 장르, 발매일, 앨범 정보 가져오기

	// db에 저장
	err = saveMelonInfoToDB(ctx, db, updatedSongs)
	if err != nil {
		log.Printf("Error saving songs to DB: %v", err)
		return
	}
}

func fetchNewSongs(db *sql.DB) ([]NewSongStruct, error) {
	now := time.Now()
	year := now.Format("2006") // YYYY
	month := now.Format("01")  // MM
	url := fmt.Sprintf("https://m.tjmedia.com/tjsong/song_monthNew.asp?YY=%s&MM=%s", year, month)

	// HTTP 요청
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %v", err)
	}
	defer res.Body.Close()

	// HTML 파싱
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	var songs []NewSongStruct
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			// 첫 번째 행은 헤더이므로 스킵
			return
		}

		cols := s.Find("td")
		if cols.Length() >= 3 {
			songNumberStr := cols.Eq(0).Text()
			songNumber, err := strconv.Atoi(songNumberStr)
			if err != nil {
				log.Printf("Invalid song number: %s, skipping...", songNumberStr)
				return
			}

			song := NewSongStruct{
				SongNumber: songNumber,
				SongName:   cols.Eq(1).Text(),
				ArtistName: cols.Eq(2).Text(),
			}
			songs = append(songs, song)
		}
	})

	log.Printf("Crawled %d songs", len(songs))

	// DB에서 이번 달에 추가된 노래 가져오기
	query := fmt.Sprintf(`
		SELECT song_number 
		FROM song_info 
		WHERE MONTH(created_at) = ? AND YEAR(created_at) = ?
	`)
	rows, err := db.Query(query, month, year)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch songs from DB: %v", err)
	}
	defer rows.Close()

	// DB에 있는 노래 번호 수집
	dbSongNumbers := make(map[int]bool)
	for rows.Next() {
		var songNumber int
		if err := rows.Scan(&songNumber); err != nil {
			return nil, fmt.Errorf("failed to scan DB rows: %v", err)
		}
		dbSongNumbers[songNumber] = true
	}

	// 크롤링된 노래 중 DB에 없는 노래만 필터링
	var newSongs []NewSongStruct
	for _, song := range songs {
		if !dbSongNumbers[song.SongNumber] {
			newSongs = append(newSongs, song)
		}
	}

	log.Printf("%d new songs found that are not in the database.", len(newSongs))
	return newSongs, nil
}

func fetchMrAndLive(songs []NewSongStruct) ([]NewSongStruct, error) {
	for i, song := range songs {
		isLive, isMr, err := fetchMrAndLiveForOne(song)
		if err != nil {
			log.Printf(err.Error())
			return nil, err
		}
		songs[i].IsLive = isLive
		songs[i].IsMr = isMr
	}

	return songs, nil
}

func fetchMrAndLiveForOne(song NewSongStruct) (bool, bool, error) {
	//log.Printf("Fetching MR and Live info for song: %d - %s by %s", song.SongNumber, song.SongName, song.ArtistName)
	url := fmt.Sprintf("https://www.tjmedia.com/tjsong/song_search_list.asp?strType=16&natType=&strText=%d&strCond=1&strSize05=100", song.SongNumber)

	// HTTP 요청
	res, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to fetch MR and Live info for song %d: %v", song.SongNumber, err)
		return false, false, err
	}
	defer res.Body.Close()

	// HTML 파싱
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Printf("Failed to parse HTML for song %d: %v", song.SongNumber, err)
		return false, false, err
	}

	table := doc.Find("div#BoardType1 > table.board_type1")
	if table.Length() == 0 {
		log.Printf("No table found for song %s", song.SongNumber)
		return false, false, fmt.Errorf("no table found")
	}

	row := table.Find("tbody > tr:nth-child(2)")
	if row.Length() == 0 {
		log.Printf("No row found for song %s", song.SongNumber)
		return false, false, fmt.Errorf("no row found")
	}

	isLive := row.Find("td img[src='/images/tjsong/live_icon.png']").Length() > 0
	isMr := row.Find("td img[src='/images/tjsong/mr_icon.png']").Length() > 0

	//log.Printf("song %d - MR: %t, Live: %t", song.SongNumber, isMr, isLive)
	return isLive, isMr, nil
}

func saveToDB(ctx context.Context, db *sql.DB, songs []NewSongStruct) error {
	for _, song := range songs {
		// 새로운 레코드 생성
		newSong := mysql.SongInfo{
			SongNumber: song.SongNumber,
			SongName:   song.SongName,
			ArtistName: song.ArtistName,
			IsMR:       null.BoolFrom(song.IsMr),
			IsLive:     null.BoolFrom(song.IsLive),
		}

		// DB에 저장
		err := newSong.Insert(ctx, db, boil.Infer())
		if err != nil {
			log.Printf("Failed to insert song %s: %v", song.SongNumber, err)
		}
	}

	log.Printf("Saved %d songs to DB", len(songs))
	return nil
}

func fetchMelonSongId(songs []NewSongStruct) ([]NewSongStruct, error) {
	for i, song := range songs {

		cleanedArtistName := cleanArtistName(song.ArtistName)         // 괄호 빼고 정리된 아티스트 이름A
		innerArtistName := extractParenthesesContent(song.ArtistName) // 괄호 안의 아티스트 이름 B
		// 원래 아티스트 이름 C

		// 1. 기본적으로 제목과 정리된 아티스트 이름(A)으로 검색
		rows, err := SearchMelon(song.SongName, cleanedArtistName)

		// 2. 여전히 결과가 없으면 괄호 안의 내용(영어 이름, B)으로 검색
		if err == nil && (rows == nil || len(rows) == 0) {
			rows, err = SearchMelon(song.SongName, innerArtistName)
		}

		// 3. 여전히 결과가 없으면 원래 아티스트 이름(C)으로 다시 검색
		if err == nil && (rows == nil || len(rows) == 0) {
			rows, err = SearchMelon(song.SongName, song.ArtistName)
		}

		// 4. 최종적으로 결과가 없으면 제목만으로 검색
		if err == nil && (rows == nil || len(rows) == 0) {
			rows, err = SearchMelon(song.SongName, "")
		}

		// 검색 결과 없다면 스킵
		if rows == nil || len(rows) == 0 {
			log.Printf("No results found for song: %s by %s", song.SongName, song.ArtistName)
			continue
		}

		// 검색 결과가 있다면 best match 찾기
		bestMatch := findHighestSimilarityMatch(song.SongName, cleanedArtistName, rows)
		if bestMatch == nil {
			log.Printf("No suitable match found for song: %s by %s", song.SongName, song.ArtistName)
			continue
		}

		// best match 정보로 melon song id 업데이트
		songs[i].MelonSongId = bestMatch.MelonSongId

		log.Printf("Updated song: %s with Melon Song ID: %s", song.SongName, bestMatch.MelonSongId)

	}

	return songs, nil
}

func cleanArtistName(artistName string) string {
	re := regexp.MustCompile(`$begin:math:text$[^)]*$end:math:text$`)
	cleanedName := re.ReplaceAllString(artistName, "")
	return strings.TrimSpace(cleanedName)
}

func extractParenthesesContent(artistName string) string {
	re := regexp.MustCompile(`$begin:math:text$([^)]*)$end:math:text$`)

	matches := re.FindAllStringSubmatch(artistName, -1)

	var validContent []string
	for _, match := range matches {
		if len(match) > 1 {
			content := strings.TrimSpace(match[1])
			lowerContent := strings.ToLower(content)

			if !strings.Contains(lowerContent, "feat") && !strings.Contains(lowerContent, "featuring") {
				validContent = append(validContent, content)
			}
		}
	}

	if len(validContent) > 0 {
		return validContent[0]
	}

	return ""
}

func SearchMelon(title, artist string) ([]SearchResult, error) {
	// User-Agent 리스트
	userAgents := []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.5735.110 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36",
	}

	// 랜덤 User-Agent 선택
	randomUserAgent := userAgents[rand.Intn(len(userAgents))]

	// HTTP 클라이언트 생성
	client := &http.Client{}
	searchURL := fmt.Sprintf("https://www.melon.com/search/song/index.htm?q=%s", url.QueryEscape(title+" "+artist))

	randomSleep()

	// HTTP 요청 생성
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		log.Printf(err.Error())
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("User-Agent", randomUserAgent)

	// HTTP 요청 전송
	resp, err := client.Do(req)
	if err != nil {
		log.Printf(err.Error())
		return nil, fmt.Errorf("failed to fetch page for %s by %s: %v", title, artist, err)
	}
	defer resp.Body.Close()

	// 상태 코드 확인
	if resp.StatusCode != http.StatusOK {
		log.Printf("failed to fetch page for %s by %s: Status code %d", title, artist, resp.StatusCode)
		return nil, fmt.Errorf("failed to fetch page for %s by %s: Status code %d", title, artist, resp.StatusCode)
	}

	// HTML 파싱
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf(err.Error())
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	// 검색 결과 추출
	rows := doc.Find("#frm_defaultList > div > table > tbody > tr")
	if rows.Length() == 0 {
		log.Printf("No results found for %s by %s", title, artist)
		return nil, nil
	}

	// 상위 3개의 결과 추출
	var topResults []SearchResult
	rows.EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i >= 3 {
			return false
		}

		// 곡 제목 추출
		songNameTag := s.Find("td:nth-child(3) > div > div > a.fc_gray")
		songName := strings.TrimSpace(songNameTag.Text())

		// 곡 ID 추출
		linkElement := s.Find("td:nth-child(3) > div > div > a.btn.btn_icon_detail")
		href, exists := linkElement.Attr("href")
		var songID string
		if exists {
			match := regexp.MustCompile(`searchLog\('web_song','SONG','SO','([^']+)','(\d+)'\);`).FindStringSubmatch(href)
			if len(match) > 2 {
				songID = match[2]
			}
		}

		// 아티스트 이름 추출
		artistNameTag := s.Find("#artistName > a")
		artistName := strings.TrimSpace(artistNameTag.Text())

		// 결과 추가
		if songName != "" && artistName != "" && songID != "" {
			topResults = append(topResults, SearchResult{
				SongName:    songName,
				ArtistName:  artistName,
				MelonSongId: songID,
			})
		}

		return true
	})

	log.Printf("Search results found for %s by %s", title, artist)

	return topResults, nil
}

func findHighestSimilarityMatch(title, artist string, results []SearchResult) *SearchResult {
	validMatches := []struct {
		Similarity float64
		Result     SearchResult
	}{}

	for _, result := range results {
		resultSongName := removeBrackets(result.SongName)
		resultArtistName := removeBrackets(result.ArtistName)

		avgSimilarity := calculateSimilarity(title, artist, resultSongName, resultArtistName)
		if avgSimilarity >= 0.6 {
			validMatches = append(validMatches, struct {
				Similarity float64
				Result     SearchResult
			}{Similarity: avgSimilarity, Result: result})
		}
		log.Printf("Title: %s, Result Title: %s, Artist: %s, Result Artist: %s, Similarity: %.2f", title, resultSongName, artist, resultArtistName, avgSimilarity)
	}

	if len(validMatches) == 0 {
		return nil
	}

	bestMatch := validMatches[0]
	for _, match := range validMatches {
		if match.Similarity > bestMatch.Similarity {
			bestMatch = match
		}
	}

	return &bestMatch.Result
}

func removeSpacesIfKorean(text string) string {
	isKorean := regexp.MustCompile(`^[가-힣\s]+$`)
	if isKorean.MatchString(strings.ReplaceAll(text, " ", "")) {
		return strings.ReplaceAll(text, " ", "")
	}
	return text
}

func removeBrackets(text string) string {
	re := regexp.MustCompile(`\(.*?\)`)
	return strings.TrimSpace(re.ReplaceAllString(text, ""))
}

func calculateSimilarity(title, artist, resultTitle, resultArtist string) float64 {
	titleSimilarity := 1.0 - float64(levenshtein.ComputeDistance(strings.ToLower(title), strings.ToLower(resultTitle)))/float64(len(title))
	artistSimilarity := 1.0 - float64(levenshtein.ComputeDistance(strings.ToLower(artist), strings.ToLower(resultArtist)))/float64(len(artist))
	return (titleSimilarity + artistSimilarity) / 2.0
}

func randomSleep() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	sleepDuration := time.Duration(r.Intn(10)+3) * time.Second
	time.Sleep(sleepDuration)
}

func saveMelonInfoToDB(ctx context.Context, db *sql.DB, songs []NewSongStruct) error {
	for _, song := range songs {
		// 기존 데이터를 조회
		existingSong, err := mysql.SongInfos(qm.Where("song_number = ?", song.SongNumber)).One(ctx, db)
		if err != nil {
			log.Printf("Failed to fetch song %d from DB: %v", song.SongNumber, err)
			continue
		}
		if existingSong == nil {
			log.Printf("Song with number %d does not exist in the database.", song.SongNumber)
			continue
		}

		// 업데이트할 데이터 설정
		existingSong.MelonSongID = null.StringFrom(song.MelonSongId)
		existingSong.Album = null.StringFrom(song.Album)
		existingSong.Genre = null.StringFrom(song.Genre)
		existingSong.Year = null.IntFrom(song.ReleaseYear)

		// 업데이트 실행
		_, err = existingSong.Update(ctx, db, boil.Infer())
		if err != nil {
			log.Printf("Failed to update song %d in DB: %v", song.SongNumber, err)
			continue
		}
	}

	log.Printf("Updated %d songs to DB", len(songs))
	return nil
}
