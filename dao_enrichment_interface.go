package main

type DaoEnrichment[T any] interface {
	GetEnrichment(t T) (res T)
}
