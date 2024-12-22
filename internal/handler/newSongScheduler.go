package handler

import (
	"SingSong-Server/internal/db/mysql"
	"context"
	"database/sql"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type NewSongStruct struct {
	SongNumber int
	SongName   string
	ArtistName string
	IsMr       bool
	IsLive     bool
}

func ScheduleNewSongs(db *sql.DB) {
	ctx := context.Background()

	// 신곡 가져오기
	newSongs, err := fetchNewSongs(db)
	if err != nil {
		log.Printf(err.Error())
		return
	}

	// isMr과 isLive 가져오기
	updatedSongs, err := fetchMrAndLive(newSongs)
	if err != nil {
		log.Printf(err.Error())
		return
	}

	// db에 1차적으로 저장
	err = saveToDB(ctx, db, updatedSongs)
	if err != nil {
		log.Printf("Error saving songs to DB: %v", err)
		return
	}

	// melon song id와 앨범 이미지 크롤링해서 저장

	// 장르, 발매일, 앨범 정보 크롤링해서 저장
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
	return songs, nil
}

func fetchMrAndLive(songs []NewSongStruct) ([]NewSongStruct, error) {
	for i, song := range songs {
		isLive, isMr, err := FetchMrAndLiveForOne(song)
		if err != nil {
			log.Printf(err.Error())
			return nil, err
		}
		songs[i].IsLive = isLive
		songs[i].IsMr = isMr
	}

	return songs, nil
}

func FetchMrAndLiveForOne(song NewSongStruct) (bool, bool, error) {
	log.Printf("Fetching MR and Live info for song: %d - %s by %s", song.SongNumber, song.SongName, song.ArtistName)
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

	log.Printf("song %d - MR: %t, Live: %t", song.SongNumber, isMr, isLive)
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
