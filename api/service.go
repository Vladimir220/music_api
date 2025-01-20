package api

import (
	"encoding/json"
	"fmt"
	"log"
	"music_api/dao/caching"
	"music_api/dao/db"
	enrch "music_api/dao/enrichment"
	"music_api/models"
	"strings"
	"time"
)

type Service struct {
	dao      db.DaoDB[models.Track]
	enrch    enrch.EnrichmentChain[models.Track]
	caching  caching.DaoCaching
	debugLog *log.Logger
}

func (s Service) CreateTrack(song, group string) (code int) {
	track := models.Track{Song: song, Group_name: group}
	var success bool
	track, success = s.enrch.Execute(track)
	if !success {
		s.debugLog.Println("Обогащение не состоялось")
	}

	rowsAffected, err := s.dao.Create(track)
	if err != nil {
		code = 500
		s.debugLog.Println(err.Error())
		return
	} else if rowsAffected == 0 {
		code = 400
		s.debugLog.Println("Данные пользователя некорректны: 400")
		return
	}

	code = 201
	return
}

func (s Service) ReadAll(filter models.Track, recordStart, recordEnd int) (recordsJSON []byte, code int) {
	query := fmt.Sprintf("ReadAll:GroupName=%s:Song=%s:ReleaseDate=%s:Link=%s:Lyrics=%s:RecordStart=%d:RecordEnd=%d", filter.Group_name, filter.Song, filter.Release_date,
		filter.Link, filter.Song_lyrics, recordStart, recordEnd)
	cach, err := s.caching.Get(query)
	if err == nil {
		s.debugLog.Println("Кэш найден!")
		return []byte(cach), 200
	} else {
		s.debugLog.Printf("Кэш не найден:%v\n", err)
	}

	res, err := s.dao.Read([]string{"*"}, &filter, recordStart, recordEnd)
	if err != nil {
		code = 500
		s.debugLog.Println(err.Error())
		return
	} else if len(res) == 0 {
		code = 404
		s.debugLog.Println("Данные не найдены: 404")
		return
	}

	records := make([]map[string]string, 0)
	for _, t := range res {
		m := make(map[string]string)

		m["song"] = t.Song
		m["group_name"] = t.Group_name

		date, err := time.Parse("2006-01-02T15:04:05Z", t.Release_date)
		if err != nil {
			code = 500
			s.debugLog.Println("Неудачный парсинг даты")
			return
		}
		m["release_date"] = date.Format("2006-01-02")

		m["song_lyrics"] = t.Song_lyrics
		m["Link"] = t.Link
		records = append(records, m)
	}

	recordsJSON, err = json.Marshal(records)
	if err != nil {
		code = 500
		s.debugLog.Println(err.Error())
		return
	}

	s.caching.Set(query, string(recordsJSON), 10*time.Minute)

	code = 200
	return
}

func (s Service) ReadTrackLyrics(song, group string, coupletStart, coupletEnd int) (lyricsJSON []byte, code int) {
	query := fmt.Sprintf("ReadTrackLyrics:GroupName=%s:Song=%s:CoupletStart=%d:CoupletEnd=%d", group, song, coupletStart, coupletEnd)
	cach, err := s.caching.Get(query)
	if err == nil {
		s.debugLog.Println("Кэш найден!")
		return []byte(cach), 200
	} else {
		s.debugLog.Printf("Кэш не найден:%v\n", err)
	}

	track := models.Track{Song: song, Group_name: group}
	res, err := s.dao.Read([]string{"Song_lyrics"}, &track, 0, 1)
	if err != nil {
		code = 500
		s.debugLog.Println(err.Error())
		return
	} else if len(res) == 0 {
		code = 404
		s.debugLog.Println("Данные не найдены: 404")
		return
	}

	cleanedLyrics := strings.ReplaceAll(res[0].Song_lyrics, "\\n", "\n")
	lyricCouplets := strings.Split(cleanedLyrics, "\n\n")

	if coupletEnd > (len(lyricCouplets)) {
		coupletEnd = len(lyricCouplets)
	}
	requiredCouplets := strings.Join(lyricCouplets[coupletStart:coupletEnd], "\n\n")

	lyrics := map[string]string{"Song_lyrics": requiredCouplets}

	lyricsJSON, err = json.Marshal(lyrics)
	if err != nil {
		code = 500
		s.debugLog.Println(err.Error())
		return
	}

	s.caching.Set(query, string(lyricsJSON), 10*time.Minute)
	code = 200
	return
}

func (s Service) DeleteTrack(song, group string) (code int) {
	t := models.Track{Song: song, Group_name: group}
	rowsAffected, err := s.dao.Delete(t)
	if err != nil {
		code = 500
		s.debugLog.Println(err.Error())
		return
	} else if rowsAffected == 0 {
		code = 400
		s.debugLog.Println("Данные пользователя некорректны: 400")
		return
	}

	code = 204
	return
}

func (s Service) UpdateTrack(song, group string, newData models.Track) (code int) {
	t := models.Track{Song: song, Group_name: group}
	rowsAffected, err := s.dao.Update(t, newData)
	if err != nil {
		code = 500
		s.debugLog.Println(err.Error())
		return
	} else if rowsAffected == 0 {
		code = 400
		s.debugLog.Println("Данные пользователя некорректны: 400")
		return
	}

	code = 204
	return
}

/*
type SongDetail struct {
	ReleaseDate string `json:"releaseDate" redis:"releaseDate"`
	Text        string `json:"text" redis:"text"`
	Link        string `json:"link" redis:"link"`
}
*/

func (s Service) ReadInfo(filter models.Track) (info models.SongDetail, code int) {
	query := fmt.Sprintf("ReadAll:GroupName=%s:Song=%s:ReleaseDate=%s:Link=%s:Lyrics=%s", filter.Group_name, filter.Song, filter.Release_date, filter.Link, filter.Song_lyrics)
	cach, err := s.caching.HGet(query)
	if err == nil {
		s.debugLog.Println("Кэш найден!")
		fmt.Println(cach)
		cachObj := models.SongDetail{ReleaseDate: cach["releaseDate"], Text: cach["text"], Link: cach["link"]}
		return cachObj, 200
	} else {
		s.debugLog.Printf("Кэш не найден:%v\n", err)
	}

	res, err := s.dao.Read([]string{"*"}, &filter, 0, 1)
	if err != nil {
		code = 500
		s.debugLog.Println(err.Error())
		return
	} else if len(res) == 0 {
		code = 404
		s.debugLog.Println("Данные не найдены: 404")
		return
	}

	t := res[0]

	info = models.SongDetail{ReleaseDate: t.Release_date, Text: t.Song_lyrics, Link: t.Link}
	var date time.Time
	date, err = time.Parse("2006-01-02T15:04:05Z", info.ReleaseDate)
	if err != nil {
		code = 500
		s.debugLog.Println("Неудачный парсинг даты")
		return
	}
	info.ReleaseDate = date.Format("2006-01-02")

	s.caching.HSet(query, info)

	code = 200
	return
}

func CreateService(dao db.DaoDB[models.Track], enrch enrch.EnrichmentChain[models.Track], caching caching.DaoCaching, debugLog *log.Logger) MusicService {
	return Service{dao: dao, enrch: enrch, debugLog: debugLog, caching: caching}
}
