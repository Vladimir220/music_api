package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type daoLastFm struct {
	url      string
	nextStep EnrichmentChain[Track]
	Input    struct {
		Body struct {
			Link string `json:"url"`
			Wiki struct {
				ReleaseDate string `json:"published"`
			} `json:"wiki"`
		} `json:"track"`
	}
}

func (d *daoLastFm) GetEnrichment(t Track) (res Track, err error) {
	token := os.Getenv("TOKEN_LASTFM")
	if token == "" {
		err = fmt.Errorf("в окружении нет токена для API LastFm: %s", token)
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

	res.Link = d.Input.Body.Link
	res.Release_date = d.Input.Body.Wiki.ReleaseDate

	return
}

func (e *daoLastFm) SetNext(next EnrichmentChain[Track]) {
	e.nextStep = next
}

func (e *daoLastFm) Execute(t Track) (res Track, success bool) {
	buf, err := e.GetEnrichment(t)
	if err != nil && e.nextStep != nil {
		buf, success = e.nextStep.Execute(t)
	}
	return buf, success
}

func createDaoLastFm() *daoLastFm {
	return &daoLastFm{url: "https://ws.audioscrobbler.com/2.0"}
}
