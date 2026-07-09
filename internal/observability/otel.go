// Package observability wires up structured logging and OpenTelemetry
// (traces + metrics) for the service. The exporter endpoint is read from the
// standard OTEL_EXPORTER_OTLP_ENDPOINT env var so that in Kubernetes it can
// simply point at "localhost:4317" (a sidecar collector injected by the
// otel-sidecar-injector controller) with zero code changes.
package observability

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// ServiceName is used as the OTel "service.name" resource attribute and as
// the annotation value the sidecar-injector controller matches on.
const ServiceName = "shorturl"

// Shutdown stops all providers created by Setup. Callers should defer it and
// call it with a bounded context during graceful shutdown.
type Shutdown func(ctx context.Context) error

// Setup configures a global slog logger (JSON) and OTel tracer/meter
// providers exporting over OTLP/gRPC. It never fails hard on missing
// collector connectivity: exporters retry/drop in the background rather
// than blocking application startup.
func Setup(ctx context.Context) (Shutdown, error) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(ServiceName),
			semconv.ServiceVersion(getEnv("APP_VERSION", "dev")),
			semconv.DeploymentEnvironment(getEnv("environment", "dev")),
		),
		resource.WithHost(),
		resource.WithProcess(),
	)
	if err != nil {
		return nil, err
	}

	traceExp, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExp),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{},
	))

	metricExp, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	mp := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(metricExp, metric.WithInterval(15*time.Second))),
		metric.WithResource(res),
	)
	otel.SetMeterProvider(mp)

	slog.Info("observability initialized",
		"otel_endpoint", getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317"),
		"service", ServiceName,
	)

	return func(shutdownCtx context.Context) error {
		var errs []error
		if err := tp.Shutdown(shutdownCtx); err != nil {
			errs = append(errs, err)
		}
		if err := mp.Shutdown(shutdownCtx); err != nil {
			errs = append(errs, err)
		}
		return errors.Join(errs...)
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
