package view

type PaginatedView[T any] struct {
	Items    []T   `json:"items"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
}

func NewPaginatedView[T any](items []T, total int64, page int, pageSize int) PaginatedView[T] {
	return PaginatedView[T]{Items: items, Total: total, Page: page, PageSize: pageSize}
}
