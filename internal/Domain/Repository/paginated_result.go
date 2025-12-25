package repository

type PaginatedResult[T any] struct {
	Items    []T
	Total    int64
	Page     int
	PageSize int
}
