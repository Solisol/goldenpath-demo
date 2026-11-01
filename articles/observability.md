---
name: Logs, metrics and traces
keywords: [observability, logging, logs, log, metrics, tracing, traces, opentelemetry, otel, monitoring, slog, health]
status: draft
---
- **Logs**: `log/slog` with the JSON handler to stdout. No `fmt.Println` for logging.
- **Traces and metrics**: OpenTelemetry, exported via OTLP to the collector at
  `OTEL_EXPORTER_OTLP_ENDPOINT` (set by the platform). For prototypes, logs are enough.
- **Health**: `GET /healthz` returns 200 when the service can serve traffic.
