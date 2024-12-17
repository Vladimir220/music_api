package main

type trackEnricherDefault struct {
	nextStep EnrichmentChain[Track]
}

func (e *trackEnricherDefault) GetEnrichment(t Track) (res Track, err error) {
	t.Link = "Enrichment has occurred"
	t.Release_date = "2011-11-11"
	t.Song_lyrics = "Enrichment has occurred"
	res = t
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
