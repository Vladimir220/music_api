package enrichment

import (
	"music_api/models"
)

type trackEnricherDefault struct {
	nextStep EnrichmentChain[models.Track]
}

func (e *trackEnricherDefault) GetEnrichment(t models.Track) (res models.Track, err error) {
	res = t
	if res.Link == "" {
		res.Link = "Enrichment has occurred"
	}
	if res.Release_date == "" {
		res.Release_date = "2011-11-11"
	}
	if res.Song_lyrics == "" {
		res.Song_lyrics = "Enrichment has occurred"
	}

	return
}

func (e *trackEnricherDefault) SetNext(next EnrichmentChain[models.Track]) {
	e.nextStep = next
}

func (e *trackEnricherDefault) Execute(t models.Track) (res models.Track, success bool) {
	res, _ = e.GetEnrichment(t)
	success = true
	return
}

func CreateTrackEnricherDefault() *trackEnricherDefault {
	return &trackEnricherDefault{}
}
