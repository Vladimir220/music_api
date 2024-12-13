package main

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
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type TrackIdentifier struct {
	Song       string `json:"song"`
	Group_name string `json:"group_name"`
}
