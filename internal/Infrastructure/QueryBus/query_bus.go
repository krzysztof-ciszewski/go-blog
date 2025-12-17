package query_bus

import (
	"context"
	"errors"
)

type QueryBus interface {
	Execute(ctx context.Context, query any) (any, error)
	RegisterHandler(handler QueryHandler)
}

type queryBus struct {
	handlers []QueryHandler
}

func NewQueryBus() QueryBus {
	return &queryBus{
		handlers: []QueryHandler{},
	}
}

func (q *queryBus) Execute(ctx context.Context, query any) (any, error) {
	for _, handler := range q.handlers {
		if handler.Supports(query) {
			return handler.Handle(ctx, query)
		}
	}

	return nil, errors.New("no handler found for query")
}

func (q *queryBus) RegisterHandler(handler QueryHandler) {
	q.handlers = append(q.handlers, handler)
}
