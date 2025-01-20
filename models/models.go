package models

type Track struct {
	Song         string `json:"song"`
	Group_name   string `json:"group_name"`
	Release_date string `json:"release_date"`
	Song_lyrics  string `json:"song_lyrics"`
	Link         string `json:"link"`
}

type Lyrics struct {
	Text string `json:"Lyrics"`
}

type SongDetail struct {
	ReleaseDate string `json:"releaseDate" redis:"releaseDate"`
	Text        string `json:"text" redis:"text"`
	Link        string `json:"link" redis:"link"`
}

type TrackIdentifier struct {
	Song       string `json:"song"`
	Group_name string `json:"group_name"`
}
