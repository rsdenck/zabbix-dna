package observability

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.39.0"
)

type OTLPEngine struct {
	Endpoint    string
	ServiceName string
}

func NewOTLPEngine(endpoint, serviceName string) *OTLPEngine {
	return &OTLPEngine{
		Endpoint:    endpoint,
		ServiceName: serviceName,
	}
}

func (e *OTLPEngine) InitMetrics(ctx context.Context) (*metric.MeterProvider, error) {
	// Strip protocol if present for WithEndpoint
	endpoint := e.Endpoint
	if len(endpoint) > 7 && endpoint[:7] == "http://" {
		endpoint = endpoint[7:]
	} else if len(endpoint) > 8 && endpoint[:8] == "https://" {
		endpoint = endpoint[8:]
	}

	exporter, err := otlpmetrichttp.New(ctx,
		otlpmetrichttp.WithEndpoint(endpoint),
		otlpmetrichttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(e.ServiceName),
		),
	)
	if err != nil {
		return nil, err
	}

	mp := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(exporter, metric.WithInterval(60*time.Second))),
		metric.WithResource(res),
	)
	otel.SetMeterProvider(mp)

	return mp, nil
}

func (e *OTLPEngine) InitTraces(ctx context.Context) (*trace.TracerProvider, error) {
	// Strip protocol if present for WithEndpoint
	endpoint := e.Endpoint
	if len(endpoint) > 7 && endpoint[:7] == "http://" {
		endpoint = endpoint[7:]
	} else if len(endpoint) > 8 && endpoint[:8] == "https://" {
		endpoint = endpoint[8:]
	}

	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(e.ServiceName),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)
	otel.SetTracerProvider(tp)

	return tp, nil
}
