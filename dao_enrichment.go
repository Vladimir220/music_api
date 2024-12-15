package main

type DaoEnrichment[T any] interface {
	GetEnrichment(t T) (res T)
}

type CreateDaoEnrichment[T any] func() DaoEnrichment[T]
