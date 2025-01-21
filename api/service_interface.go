package api

import (
	"log"

	"music_api/dao/caching"
	"music_api/dao/db"
	enrc "music_api/dao/enrichment"
	"music_api/models"
)

type MusicService interface {
	CreateTrack(song, group string) (code int)
	ReadAll(filter models.Track, recordStart, recordEnd int) (recordsJSON []byte, code int)
	ReadTrackLyrics(song, group string, coupletStart, coupletEnd int) (lyricsJSON []byte, code int)
	DeleteTrack(song, group string) (code int)
	UpdateTrack(song, group string, newData models.Track) (code int)
}

type CreateMusicService func(dao db.DaoDB[models.Track], enrch enrc.EnrichmentChain[models.Track], caching caching.DaoCaching, debugLog *log.Logger) MusicService
