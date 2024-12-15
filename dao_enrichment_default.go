package main

type trackEnricherDefault struct {
}

func (e *trackEnricherDefault) GetEnrichment(t Track) (res Track) {
	t.Link = "Enrichment has occurred"
	t.Release_date = "2011-11-11"
	t.Song_lyrics = "Enrichment has occurred"
	res = t
	return
}

func createTrackEnricherDefault() *trackEnricherDefault {
	return &trackEnricherDefault{}
}
