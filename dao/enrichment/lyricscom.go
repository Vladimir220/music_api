package enrichment

import (
	"encoding/json"
	"fmt"
	"io"
	"music_api/models"
	"net/http"
	"net/url"
	"os"
)

type daoLyricsCom struct {
	url      string
	nextStep EnrichmentChain[models.Track]
	Input    struct {
		Body []struct {
			SongLyrics string `json:"song-link"`
		} `json:"result"`
	}
}

func (d *daoLyricsCom) GetEnrichment(t models.Track) (res models.Track, err error) {
	token := os.Getenv("TOKEN_LYRICSCOM")
	uid := os.Getenv("UID_LYRICSCOM")
	if token == "" || uid == "" {
		err = fmt.Errorf("в окружении нет токена или uid для API LyricsCom")
		return
	}

	res = t

	fullUrl, _ := url.Parse(d.url)
	params := url.Values{}
	params.Add("uid", uid)
	params.Add("tokenid", token)
	params.Add("term", t.Song)
	params.Add("artist", t.Group_name)
	params.Add("format", "json")
	fullUrl.RawQuery = params.Encode()

	var resp *http.Response
	resp, err = http.Get(fullUrl.String())
	if err != nil {
		err = fmt.Errorf("не удалось получить доступ к ресурсу обогащения: %v", err)
		return
	}

	var body []byte
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("не удалось получить доступ к ресурсу обогащения: %v", err)
		return
	}
	defer resp.Body.Close()

	err = json.Unmarshal(body, &d.Input)
	if err != nil {
		err = fmt.Errorf("не удалось распоковать ответ ресурса обогащения: %v", err)
		return
	}

	if len(d.Input.Body) == 0 {
		err = fmt.Errorf("текст песни не найден")
		return
	}

	res.Song_lyrics = d.Input.Body[0].SongLyrics

	return
}

func (e *daoLyricsCom) SetNext(next EnrichmentChain[models.Track]) {
	e.nextStep = next
}

func (e *daoLyricsCom) Execute(t models.Track) (res models.Track, success bool) {
	buf, err := e.GetEnrichment(t)
	condContin := (err != nil || buf.Link == "" || buf.Release_date == "" || buf.Song_lyrics == "") && e.nextStep != nil
	if condContin {
		buf, success = e.nextStep.Execute(buf)
	}
	return buf, success
}

func CreateDaoLyricsCom() *daoLyricsCom {
	return &daoLyricsCom{url: "https://www.stands4.com/services/v2/lyrics.php"}
}
