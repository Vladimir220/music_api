package enrichment

import (
	"encoding/json"
	"fmt"
	"music_api/models"
	"net/http"
	"net/url"
	"os"
)

type enrchServer struct {
	nextStep EnrichmentChain[models.Track]
}

func (d *enrchServer) GetEnrichment(t models.Track) (res models.Track, err error) {
	host := os.Getenv("ENRCH_SERVER_HOST")
	if host == "" {
		err = fmt.Errorf("в окружении нет адреса хоста стандартного обогатителя")
		return
	}

	res = t

	fullUrl, _ := url.Parse(host)
	params := url.Values{}
	params.Add("group", t.Group_name)
	params.Add("song", t.Song)
	fullUrl.RawQuery = params.Encode()

	resp, err := http.Get(fullUrl.String())
	if err != nil {
		err = fmt.Errorf("не удалось получить доступ к ресурсу обогащения: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("информация не найдена: %v", err)
		return
	}

	var sd models.SongDetail
	err = json.NewDecoder(resp.Body).Decode(&sd)
	if err != nil {
		err = fmt.Errorf("не удалось распоковать ответ ресурса обогащения: %v", err)
		return
	}

	res.Link = sd.Link
	res.Release_date = sd.ReleaseDate
	res.Song_lyrics = sd.Text

	return
}

func (e *enrchServer) SetNext(next EnrichmentChain[models.Track]) {
	e.nextStep = next
}

func (e *enrchServer) Execute(t models.Track) (res models.Track, success bool) {
	buf, err := e.GetEnrichment(t)
	condContin := (err != nil || buf.Link == "" || buf.Release_date == "" || buf.Song_lyrics == "")
	if condContin && e.nextStep != nil {
		buf, success = e.nextStep.Execute(buf)
		return
	} else if condContin {
		return
	}
	return buf, true
}

func CreateDaoEnrchServer() *daoLastFm {
	return &daoLastFm{}
}
