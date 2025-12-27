package query_bus

import (
	"context"
	"errors"
	open_telemetry "main/internal/Infrastructure/OpenTelemetry"
	"reflect"
)

type QueryBus interface {
	Execute(ctx context.Context, query any) (any, error)
	RegisterHandler(handler QueryHandler)
}

type queryBus struct {
	handlers  []QueryHandler
	telemetry open_telemetry.Telemetry
}

func NewQueryBus(telemetry open_telemetry.Telemetry) QueryBus {
	return &queryBus{
		handlers:  []QueryHandler{},
		telemetry: telemetry,
	}
}

func (q *queryBus) Execute(ctx context.Context, query any) (any, error) {
	_, span := q.telemetry.TraceStart(ctx, "QueryBus.Execute."+reflect.TypeOf(query).Name())
	defer span.End()

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
