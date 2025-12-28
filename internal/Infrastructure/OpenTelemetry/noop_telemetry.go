package open_telemetry

import (
	"context"
	config "main/internal/Infrastructure/Config"
	"os"

	"github.com/gin-gonic/gin"
	otelmetric "go.opentelemetry.io/otel/metric"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type NoopTelemetry struct {
	serviceName string
}

func NewNoopTelemetry(config config.TelemetryConfig) *NoopTelemetry {
	return &NoopTelemetry{serviceName: config.ServiceName}
}

func (n *NoopTelemetry) GetServiceName() string {
	return n.serviceName
}

func (n *NoopTelemetry) LogInfo(args ...any) {
}

func (n *NoopTelemetry) LogErrorln(args ...any) {
}

func (n *NoopTelemetry) LogFatalln(args ...any) {
	os.Exit(1)
}

func (n *NoopTelemetry) MeterInt64Histogram(metric Metric) (otelmetric.Int64Histogram, error) {
	return nil, nil
}

func (n *NoopTelemetry) MeterInt64UpDownCounter(metric Metric) (otelmetric.Int64UpDownCounter, error) {
	return nil, nil
}

func (n *NoopTelemetry) TraceStart(ctx context.Context, name string) (context.Context, oteltrace.Span) {
	return ctx, oteltrace.SpanFromContext(ctx)
}

func (n *NoopTelemetry) LogRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func (n *NoopTelemetry) MeterRequestDuration() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func (n *NoopTelemetry) MeterRequestsInFlught() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func (n *NoopTelemetry) Shutdown(ctx context.Context) {
}
