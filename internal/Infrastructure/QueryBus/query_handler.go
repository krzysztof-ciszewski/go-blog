package query_bus

import "context"

type QueryHandler interface {
	Handle(ctx context.Context, query any) (any, error)
	Supports(query any) bool
}