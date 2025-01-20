package enrichment

type DaoEnrichment[T any] interface {
	GetEnrichment(t T) (res T, err error)
}

type EnrichmentChain[T any] interface {
	SetNext(EnrichmentChain[T])
	Execute(t T) (res T, success bool)
}
