package main

type trackEnricherDefault struct {
	nextStep EnrichmentChain[Track]
}

func (e *trackEnricherDefault) GetEnrichment(t Track) (res Track, err error) {
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

func (e *trackEnricherDefault) SetNext(next EnrichmentChain[Track]) {
	e.nextStep = next
}

func (e *trackEnricherDefault) Execute(t Track) (res Track, success bool) {
	res, _ = e.GetEnrichment(t)
	success = true
	return
}

func createTrackEnricherDefault() *trackEnricherDefault {
	return &trackEnricherDefault{}
}
