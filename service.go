package main

import (
	"encoding/json"
	"log"
	"strings"
	"time"
)

type service struct {
	dao      DaoDB[Track]
	enrch    DaoEnrichment[Track]
	debugLog *log.Logger
}

func (s service) CreateTrack(song, group string) (code int) {
	track := Track{Song: song, Group_name: group}

	track = s.enrch.GetEnrichment(track)

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

func (s service) ReadAll(filter Track, recordStart, recordEnd int) (recordsJSON []byte, code int) {
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

	code = 200
	return
}

func (s service) ReadTrackLyrics(song, group string, coupletStart, coupletEnd int) (lyricsJSON []byte, code int) {
	track := Track{Song: song, Group_name: group}
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

	code = 200
	return
}

func (s service) DeleteTrack(song, group string) (code int) {
	t := Track{Song: song, Group_name: group}
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

func (s service) UpdateTrack(song, group string, newData Track) (code int) {
	t := Track{Song: song, Group_name: group}
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

func (s service) ReadInfo(filter Track) (info SongDetail, code int) {
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

	info = SongDetail{ReleaseDate: t.Release_date, Text: t.Song_lyrics, Link: t.Link}
	var date time.Time
	date, err = time.Parse("2006-01-02T15:04:05Z", info.ReleaseDate)
	if err != nil {
		code = 500
		s.debugLog.Println("Неудачный парсинг даты")
		return
	}
	info.ReleaseDate = date.Format("2006-01-02")

	code = 200
	return
}

func createService(dao DaoDB[Track], enrch DaoEnrichment[Track], debugLog *log.Logger) MusicService {
	return service{dao: dao, enrch: enrch, debugLog: debugLog}
}
