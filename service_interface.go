package main

import "log"

type MusicService interface {
	CreateTrack(song, group string) (code int)
	ReadAll(filter Track, recordStart, recordEnd int) (recordsJSON []byte, code int)
	ReadTrackLyrics(song, group string, coupletStart, coupletEnd int) (lyricsJSON []byte, code int)
	DeleteTrack(song, group string) (code int)
	UpdateTrack(song, group string, newData Track) (code int)
	ReadInfo(filter Track) (info SongDetail, code int)
}

type CreateMusicService func(dao DaoDB[Track], enrch EnrichmentChain[Track], debugLog *log.Logger) MusicService
