package api

import (
	"encoding/json"
	"fmt"
	"log"
	"music_api/dao/caching"
	"music_api/dao/db"
	enrch "music_api/dao/enrichment"
	"music_api/models"

	"net/http"
	"strconv"
)

type handlers struct {
	daoDB           db.DaoDB[models.Track]
	daoEnrch        enrch.EnrichmentChain[models.Track]
	caching         caching.DaoCaching
	srvcConstructor CreateMusicService
	infoLog         *log.Logger
	debugLog        *log.Logger
}

// @Summary Get all Tracks
// @Description Get a list of all Tracks with optional strict filter for songs, groups, releases, lyrics, and links
// @ID get-all-tracks
// @Accept  json
// @Produce  json
// @Param start query int true "Start index"
// @Param end query int true "End index"
// @Param song-filter query string false "Strict filter by song"
// @Param group-filter query string false "Strict filter by group"
// @Param release-filter query string false "Strict filter by release"
// @Param lyrics-filter query string false "Strict filter by lyrics"
// @Param link-filter query string false "Strict filter by link"
// @Success 200 {array} Track "Successfully retrieved Track"
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /all/pages [get]
func (h handlers) getAll(w http.ResponseWriter, r *http.Request) {
	s := h.srvcConstructor(h.daoDB, h.daoEnrch, h.caching, h.debugLog)
	w.Header().Set("Content-Type", "application/json")

	h.debugLog.Println("Получение всех записей в заданном диапазоне...")

	// ----input----
	start := 0
	startTxt := r.URL.Query().Get("start")
	if startTxt != "" {
		var err error
		start, err = strconv.Atoi(startTxt)
		if err != nil {
			h.debugLog.Println("Недопустимое значение параметра start=", startTxt)
			w.WriteHeader(400)
			fmt.Fprint(w, "Error 400")
			return
		}
	} else {
		h.debugLog.Println("Не указан параметр start")
		w.WriteHeader(400)
		fmt.Fprint(w, "Error 400")
		return
	}

	end := 0
	endTxt := r.URL.Query().Get("end")
	if endTxt != "" {
		var err error
		end, err = strconv.Atoi(endTxt)
		if err != nil {
			h.debugLog.Println("Недопустимое значение параметра end=", endTxt)
			w.WriteHeader(400)
			fmt.Fprint(w, "Error 400")
			return
		}
	} else {
		h.debugLog.Println("Не указан параметр end")
		w.WriteHeader(400)
		fmt.Fprint(w, "Error 400")
		return
	}

	filter := models.Track{}

	songFilterTxt := r.URL.Query().Get("song-filter")
	if songFilterTxt != "" {
		filter.Song = songFilterTxt
	}

	groupFilterTxt := r.URL.Query().Get("group-filter")
	if groupFilterTxt != "" {
		filter.Group_name = groupFilterTxt
	}

	releaseFilterTxt := r.URL.Query().Get("release-filter")
	if releaseFilterTxt != "" {
		filter.Release_date = releaseFilterTxt
	}

	lyricsFilterTxt := r.URL.Query().Get("lyrics-filter")
	if lyricsFilterTxt != "" {
		filter.Song_lyrics = lyricsFilterTxt
	}

	linkFilterTxt := r.URL.Query().Get("link-filter")
	if linkFilterTxt != "" {
		filter.Link = linkFilterTxt
	}

	h.debugLog.Println("start=", startTxt, ", end=", endTxt, ", filter=", filter)

	// ----processing----
	res, code := s.ReadAll(filter, start, end)
	w.WriteHeader(code)

	if code != 200 {
		fmt.Fprint(w, "Error code:", code)
	} else {
		h.infoLog.Println("Получение всех записей в диапазоне от", start, "до", end, "с фильтрацией", filter)
		w.Write(res)
	}
}

// @Summary Get Track Lyrics
// @Description Get Track Lyrics
// @ID get-track-lyrics
// @Accept  json
// @Produce  json
// @Param song query string true "song"
// @Param group query string true "group"
// @Param start query int true "Start index"
// @Param end query int true "End index"
// @Success 200 {object} Lyrics "Success"
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /track/lyrics/couplets [get]
func (h handlers) getTrackLyrics(w http.ResponseWriter, r *http.Request) {
	s := h.srvcConstructor(h.daoDB, h.daoEnrch, h.caching, h.debugLog)
	w.Header().Set("Content-Type", "application/json")

	h.debugLog.Println("Получение куплетов текста песни в заданном диапазоне...")

	// ----input----
	song := r.URL.Query().Get("song")
	if song == "" {
		h.debugLog.Println("Не указан параметр song")
		w.WriteHeader(400)
		fmt.Fprint(w, "Error 400")
		return
	}

	group := r.URL.Query().Get("group")
	if group == "" {
		h.debugLog.Println("Не указан параметр group")
		w.WriteHeader(400)
		fmt.Fprint(w, "Error 400")
		return
	}

	start := 0
	startTxt := r.URL.Query().Get("start")
	if startTxt != "" {
		var err error
		start, err = strconv.Atoi(startTxt)
		if err != nil {
			h.debugLog.Println("Недопустимое значение параметра start=", startTxt)
			w.WriteHeader(400)
			fmt.Fprint(w, "Error 400")
			return
		}
	} else {
		h.debugLog.Println("Не указан параметр start")
		w.WriteHeader(400)
		fmt.Fprint(w, "Error 400")
		return
	}

	end := 0
	endTxt := r.URL.Query().Get("end")
	if endTxt != "" {
		var err error
		end, err = strconv.Atoi(endTxt)
		if err != nil {
			h.debugLog.Println("Недопустимое значение параметра end=", endTxt)
			w.WriteHeader(400)
			fmt.Fprint(w, "Error 400")
			return
		}
	} else {
		h.debugLog.Println("Не указан параметр end")
		w.WriteHeader(400)
		fmt.Fprint(w, "Error 400")
		return
	}

	h.debugLog.Println("song=", song, ", group=", group, ", start=", startTxt, ", end=", endTxt)

	// ----processing----
	res, code := s.ReadTrackLyrics(song, group, start, end)

	w.WriteHeader(code)

	if code != 200 {
		fmt.Fprint(w, "Error code:", code)
	} else {
		h.infoLog.Println("Получение куплетов текста песни", song, "группы", group, "в диапазоне от", startTxt, "до", endTxt)
		w.Write(res)
	}
}

// @Summary Delete Track
// @Description Delete Track
// @ID delete-track
// @Accept  json
// @Produce  json
// @Param delTrack body TrackIdentifier true "Delete Track"
// @Success 201
// @Failure 400
// @Failure 500
// @Router /track [delete]
func (h handlers) deleteTrack(w http.ResponseWriter, r *http.Request) {
	s := h.srvcConstructor(h.daoDB, h.daoEnrch, h.caching, h.debugLog)
	w.Header().Set("Content-Type", "application/json")

	h.debugLog.Println("Удаление песни...")

	// ----input----
	var body models.TrackIdentifier
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		h.debugLog.Println("Полученное тело не в JSON")
		w.WriteHeader(400)
		fmt.Fprint(w, "Error 400")
	}
	defer r.Body.Close()

	h.debugLog.Println("song=", body.Song, ", group=", body.Group_name)

	// ----processing----
	code := s.DeleteTrack(body.Song, body.Group_name)

	w.WriteHeader(code)

	if code != 201 {
		fmt.Fprint(w, "Error code:", code)
	} else {
		fmt.Fprint(w, "Success:", code)
		h.infoLog.Println("Удаление песни", body.Song, "группы", body.Group_name)
	}
}

// @Summary Update Track
// @Description Update Track
// @ID update-track
// @Accept  json
// @Produce  json
// @Param song query string true "song"
// @Param group query string true "group"
// @Param newvalues body Track true "new values"
// @Success 204
// @Failure 400
// @Failure 500
// @Router /track [PATCH]
func (h handlers) updateTrack(w http.ResponseWriter, r *http.Request) {
	s := h.srvcConstructor(h.daoDB, h.daoEnrch, h.caching, h.debugLog)
	w.Header().Set("Content-Type", "application/json")

	h.debugLog.Println("Изменение информации о песне...")

	// ----input----
	song := r.URL.Query().Get("song")
	if song == "" {
		h.debugLog.Println("Не указан параметр song")
		w.WriteHeader(400)
		fmt.Fprint(w, "Error 400")
		return
	}

	group := r.URL.Query().Get("group")
	if group == "" {
		h.debugLog.Println("Не указан параметр group")
		w.WriteHeader(400)
		fmt.Fprint(w, "Error 400")
		return
	}

	var body models.Track
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		h.debugLog.Println("Полученное тело не в JSON")
		w.WriteHeader(400)
		fmt.Fprint(w, "Error 400")
	}
	defer r.Body.Close()

	h.debugLog.Println("song=", song, ", group=", group, ", newValues=", body)

	// ----processing----
	code := s.UpdateTrack(song, group, body)

	w.WriteHeader(code)

	if code != 204 {
		fmt.Fprint(w, "Error code:", code)
	} else {
		fmt.Fprint(w, "Success:", code)
		h.infoLog.Println("Редактирование данных о песне", song, "группы", group, ". Новые значения:", body)
	}
}

// @Summary Create Track
// @Description Create Track
// @ID create-track
// @Accept  json
// @Produce  json
// @Param newTrack body TrackIdentifier true "new values"
// @Success 201
// @Failure 400
// @Failure 500
// @Router /track [post]
func (h handlers) createTrack(w http.ResponseWriter, r *http.Request) {
	s := h.srvcConstructor(h.daoDB, h.daoEnrch, h.caching, h.debugLog)
	w.Header().Set("Content-Type", "application/json")

	h.debugLog.Println("Добавление новой песни...")

	// ----input----
	var body models.TrackIdentifier
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		h.debugLog.Println("Полученное тело не в JSON")
		w.WriteHeader(400)
		fmt.Fprint(w, "Error 400")
	}
	defer r.Body.Close()

	h.debugLog.Println("song=", body.Song, ", group=", body.Group_name)

	// ----processing----
	code := s.CreateTrack(body.Song, body.Group_name)

	w.WriteHeader(code)

	if code != 201 {
		fmt.Fprint(w, "Error code:", code)
	} else {
		fmt.Fprint(w, "Success:", code)
		h.infoLog.Println("Добавление песни", body.Song, "группы", body.Group_name)
	}
}

// @Summary Get Info
// @Description Get Info
// @ID get-info
// @Accept  json
// @Produce  json
// @Param song query string true "song"
// @Param group query string true "group"
// @Success 200 {object} SongDetail "Success"
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /info [get]
func (h handlers) getInfo(w http.ResponseWriter, r *http.Request) {
	s := h.srvcConstructor(h.daoDB, h.daoEnrch, h.caching, h.debugLog)
	w.Header().Set("Content-Type", "application/json")

	h.debugLog.Println("Получение информации о песне...")

	// ----input----
	song := r.URL.Query().Get("song")
	if song == "" {
		h.debugLog.Println("Не указан параметр song")
		w.WriteHeader(400)
		fmt.Fprint(w, "Error 400")
		return
	}

	group := r.URL.Query().Get("group")
	if group == "" {
		h.debugLog.Println("Не указан параметр group")
		w.WriteHeader(400)
		fmt.Fprint(w, "Error 400")
		return
	}

	h.debugLog.Println("song=", song, ", group=", group)

	// ----processing----
	filter := models.Track{Song: song, Group_name: group}
	res, code := s.ReadInfo(filter)

	infoJSON, err := json.Marshal(res)
	if err != nil {
		h.debugLog.Println("Неудачный парсинг в JSON")
		w.WriteHeader(500)
		fmt.Fprint(w, "Error 500")
		return
	}

	w.WriteHeader(code)

	if code != 200 {
		fmt.Fprint(w, "Error code:", code)
	} else {
		h.infoLog.Println("Получение инфирмации о песне", song, "группы", group)
		w.Write(infoJSON)
	}
}

func CreateHandlers(daoDB db.DaoDB[models.Track], daoEnrch enrch.EnrichmentChain[models.Track], caching caching.DaoCaching, srvcConstructor CreateMusicService, infoLog *log.Logger, debugLog *log.Logger) handlers {
	return handlers{daoDB: daoDB, daoEnrch: daoEnrch, srvcConstructor: srvcConstructor, infoLog: infoLog, debugLog: debugLog, caching: caching}
}
