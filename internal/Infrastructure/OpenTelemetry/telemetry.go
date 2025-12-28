package open_telemetry

import (
	"context"
	"fmt"
	config "main/internal/Infrastructure/Config"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/bridges/otelzap"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.20.0/httpconv"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type TelemetryProvider interface {
	GetServiceName() string
	LogInfo(args ...any)
	LogErrorln(args ...any)
	LogFatalln(args ...any)
	MeterInt64Histogram(metric Metric) (otelmetric.Int64Histogram, error)
	MeterInt64UpDownCounter(metric Metric) (otelmetric.Int64UpDownCounter, error)
	TraceStart(ctx context.Context, name string) (context.Context, oteltrace.Span)
	LogRequest() gin.HandlerFunc
	MeterRequestDuration() gin.HandlerFunc
	MeterRequestsInFlught() gin.HandlerFunc
	Shutdown(ctx context.Context)
}

type Telemetry struct {
	loggerProvider *log.LoggerProvider
	metricProvider *metric.MeterProvider
	traceProvider  *trace.TracerProvider
	logger         *zap.SugaredLogger
	meter          otelmetric.Meter
	tracer         oteltrace.Tracer
	config         config.TelemetryConfig
}

func NewTelemetry(ctx context.Context, config config.TelemetryConfig) (*Telemetry, error) {
	resource := newResource(config.ServiceName, config.ServiceVersion)

	loggerProvider, err := newLoggerProvider(ctx, resource)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger provider: %w", err)
	}

	meterProvider, err := newMeterProvider(ctx, resource)
	if err != nil {
		return nil, fmt.Errorf("failed to create metric provider: %w", err)
	}
	meter := meterProvider.Meter(config.ServiceName)

	tracerProvider, err := newTracerProvider(ctx, resource)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace provider: %w", err)
	}
	tracer := tracerProvider.Tracer(config.ServiceName)

	logger := zap.New(
		zapcore.NewTee(
			zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), zapcore.AddSync(os.Stdout), zapcore.InfoLevel),
			otelzap.NewCore(config.ServiceName, otelzap.WithLoggerProvider(loggerProvider)),
		),
	)

	return &Telemetry{
		loggerProvider: loggerProvider,
		metricProvider: meterProvider,
		traceProvider:  tracerProvider,
		logger:         logger.Sugar(),
		meter:          meter,
		tracer:         tracer,
		config:         config,
	}, nil
}

func (t *Telemetry) GetServiceName() string {
	return t.config.ServiceName
}

func (t *Telemetry) LogInfo(args ...any) {
	t.logger.Info(args...)
}

func (t *Telemetry) LogErrorln(args ...any) {
	t.logger.Errorln(args...)
}

func (t *Telemetry) LogFatalln(args ...any) {
	t.logger.Fatalln(args...)
}

func (t *Telemetry) MeterInt64Histogram(metric Metric) (otelmetric.Int64Histogram, error) {
	histogram, err := t.meter.Int64Histogram(
		metric.Name,
		otelmetric.WithDescription(metric.Description),
		otelmetric.WithUnit(metric.Unit),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create int64 histogram: %w", err)
	}
	return histogram, nil
}

func (t *Telemetry) MeterInt64UpDownCounter(metric Metric) (otelmetric.Int64UpDownCounter, error) {
	counter, err := t.meter.Int64UpDownCounter(
		metric.Name,
		otelmetric.WithDescription(metric.Description),
		otelmetric.WithUnit(metric.Unit),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create int64 up down counter: %w", err)
	}
	return counter, nil
}

func (t *Telemetry) TraceStart(ctx context.Context, name string) (context.Context, oteltrace.Span) {
	return t.tracer.Start(ctx, name)
}

func (t *Telemetry) LogRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		req := c.Request
		writer := c.Writer

		remoteUser := "-"
		if u, _, ok := req.BasicAuth(); ok && u != "" {
			remoteUser = u
		}
		remoteAddr := c.ClientIP()
		if remoteAddr == "" {
			remoteAddr = "-"
		}

		timestamp := time.Now().Format("02/Jan/2006:15:04:05 -0700")

		referrer := req.Referer()
		if referrer == "" {
			referrer = "-"
		}
		userAgent := req.UserAgent()
		if userAgent == "" {
			userAgent = "-"
		}

		size := max(writer.Size(), 0)

		logLine := fmt.Sprintf("%s - %s [%s] \"%s %s %s\" %d %d \"%s\" \"%s\"",
			remoteAddr,
			remoteUser,
			timestamp,
			req.Method,
			req.RequestURI,
			req.Proto,
			writer.Status(),
			size,
			referrer,
			userAgent,
		)

		t.LogInfo(logLine)
	}
}

func (t *Telemetry) MeterRequestDuration() gin.HandlerFunc {
	histogram, err := t.MeterInt64Histogram(MetricRequestDurationMillis)
	if err != nil {
		t.LogErrorln("failed to create request duration histogram: %w", err)
	}

	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		histogram.Record(c.Request.Context(), duration.Milliseconds(), otelmetric.WithAttributes(
			httpconv.ServerRequest(t.GetServiceName(), c.Request)...,
		))
	}
}

func (t *Telemetry) MeterRequestsInFlught() gin.HandlerFunc {
	counter, err := t.MeterInt64UpDownCounter(MetricRequestsInFlight)
	if err != nil {
		t.LogErrorln("failed to create requests in flight counter: %w", err)
	}

	return func(c *gin.Context) {
		attributes := otelmetric.WithAttributes(httpconv.ServerRequest(t.GetServiceName(), c.Request)...)
		counter.Add(c.Request.Context(), 1, attributes)
		c.Next()
		counter.Add(c.Request.Context(), -1, attributes)
	}
}

func (t *Telemetry) Shutdown(ctx context.Context) {
	t.loggerProvider.Shutdown(ctx)
	t.metricProvider.Shutdown(ctx)
	t.traceProvider.Shutdown(ctx)
}
