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
	url   string
	Input struct {
		Body struct {
			Link string `json:"url"`
			Wiki struct {
				ReleaseDate string `json:"published"`
			} `json:"wiki"`
		} `json:"track"`
	}
}

func (d *daoLastFm) GetEnrichment(t Track) (res Track, err error) {

	res = t

	fullUrl, _ := url.Parse(d.url)
	params := url.Values{}
	params.Add("method", "track.getInfo")
	params.Add("api_key", os.Getenv("TOKEN_LASTFM"))
	params.Add("artist", t.Group_name)
	params.Add("track", t.Song)
	params.Add("format", "json")
	fullUrl.RawQuery = params.Encode()
	fmt.Println(fullUrl.String())

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
	//input := make(map[string]map[string]interface{})
	err = json.Unmarshal(body, &d.Input)
	if err != nil {
		err = fmt.Errorf("не удалось распоковать ответ ресурса обогащения: %v", err)
		return
	}

	res.Link = d.Input.Body.Link
	res.Release_date = d.Input.Body.Wiki.ReleaseDate
	fmt.Println(d.Input)

	return
}

func createDaoLastFm() *daoLastFm {
	return &daoLastFm{url: "https://ws.audioscrobbler.com/2.0"}
}
