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

type daoLastFm struct {
	url      string
	nextStep EnrichmentChain[models.Track]
	Input    struct {
		Body struct {
			Link string `json:"url"`
			Wiki struct {
				ReleaseDate string `json:"published"`
			} `json:"wiki"`
		} `json:"track"`
	}
}

func (d *daoLastFm) GetEnrichment(t models.Track) (res models.Track, err error) {
	token := os.Getenv("TOKEN_LASTFM")
	if token == "" {
		err = fmt.Errorf("в окружении нет токена для API LastFm")
		return
	}

	res = t

	fullUrl, _ := url.Parse(d.url)
	params := url.Values{}
	params.Add("method", "track.getInfo")
	params.Add("api_key", token)
	params.Add("artist", t.Group_name)
	params.Add("track", t.Song)
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

	buf1 := d.Input.Body.Link
	buf2 := d.Input.Body.Wiki.ReleaseDate
	if buf1 == "" || buf2 == "" {
		err = fmt.Errorf("информация о песне не найдена")
		return
	}

	res.Link = buf1
	res.Release_date = buf2

	return
}

func (e *daoLastFm) SetNext(next EnrichmentChain[models.Track]) {
	e.nextStep = next
}

func (e *daoLastFm) Execute(t models.Track) (res models.Track, success bool) {
	buf, err := e.GetEnrichment(t)
	condContin := (err != nil || buf.Link == "" || buf.Release_date == "" || buf.Song_lyrics == "") && e.nextStep != nil
	if condContin {
		buf, success = e.nextStep.Execute(buf)
	}
	return buf, success
}

func CreateDaoLastFm() *daoLastFm {
	return &daoLastFm{url: "https://ws.audioscrobbler.com/2.0"}
}
