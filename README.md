# MeshGate

**Reliability-focused HTTP API gateway written in Go.**

MeshGate is an engineering project for studying the failure modes and control mechanisms that sit between clients and unreliable or overloaded upstream services.

[![CI](https://github.com/pavankumarhfl-hub/MeshGate/actions/workflows/ci.yml/badge.svg)](https://github.com/pavankumarhfl-hub/MeshGate/actions/workflows/ci.yml)

> **Status:** Early engineering build. It is not presented as production-ready software.

## Why this project exists

A reverse proxy is easy to demonstrate. A useful gateway has to make deliberate decisions about timeouts, retries, overload, request limits, partial failure, and shutdown behavior.

MeshGate is being built around those engineering questions rather than around a long list of framework features.

## Current capabilities

- HTTP request routing to configured upstreams
- Per-process rate limiting
- Upstream request timeout
- Overall request deadline
- Maximum request-body size
- Bounded retries with jitter for retry-safe HTTP methods
- Circuit breaker with a single half-open probe
- Health and readiness endpoints
- Route inspection endpoint
- Graceful shutdown on SIGINT/SIGTERM
- Race-detector CI
- `go vet` CI
- Non-root distroless container image
- Unit and integration-style tests around failure behavior

## Architecture

```text
                    ┌──────────────────────┐
                    │       Client         │
                    └──────────┬───────────┘
                               │ HTTP
                               ▼
                    ┌──────────────────────┐
                    │      MeshGate        │
                    │                      │
                    │  deadline / limits  │
                    │  rate limiting      │
                    │  route selection    │
                    │  retry policy       │
                    │  circuit breaker    │
                    └──────────┬───────────┘
                               │
                         HTTP / upstream
                               ▼
                    ┌──────────────────────┐
                    │   Upstream service   │
                    └──────────────────────┘
```

See [`docs/architecture.md`](docs/architecture.md) for component boundaries and [`docs/reliability.md`](docs/reliability.md) for failure semantics.

## Reliability model

### Timeouts

MeshGate has separate limits for HTTP header parsing, overall request processing, and upstream communication. The goal is to prevent a slow dependency from consuming a gateway worker indefinitely.

### Retries

Retries are deliberately bounded and use jittered backoff. By default, retries are enabled only for methods commonly treated as retry-safe: `GET`, `HEAD`, `OPTIONS`, `PUT`, `DELETE`, and `TRACE`.

`POST` is not automatically retried because repeating a non-idempotent operation can duplicate side effects.

### Circuit breaker

Repeated upstream failures open the circuit. After the cooldown, one request is admitted as a half-open probe; concurrent requests remain rejected until that probe succeeds or fails.

### Resource limits

Request bodies are bounded before the proxy attempts to read them. This prevents an otherwise simple retry mechanism from turning request buffering into an unbounded memory risk.

## Configuration

MeshGate is configured through environment variables for the current single-route deployment model:

| Variable | Default | Purpose |
|---|---:|---|
| `MESHGATE_ADDR` | `:8080` | Listen address |
| `MESHGATE_REQUEST_TIMEOUT` | `10s` | Overall request deadline |
| `MESHGATE_UPSTREAM_TIMEOUT` | `3s` | Upstream client timeout |
| `MESHGATE_MAX_RETRIES` | `2` | Maximum retry attempts after the first request |
| `MESHGATE_RATE_LIMIT` | `100` | Requests per fixed one-second window per key |
| `MESHGATE_BURST` | `100` | Burst ceiling for the local limiter |
| `MESHGATE_MAX_BODY_BYTES` | `1048576` | Maximum request body size |
| `MESHGATE_ROUTE_PATH` | unset | Gateway route path, e.g. `/api` |
| `MESHGATE_ROUTE_TARGET` | unset | Upstream URL, e.g. `http://127.0.0.1:9000` |

Example:

```bash
MESHGATE_ROUTE_PATH=/api \
MESHGATE_ROUTE_TARGET=http://127.0.0.1:9000 \
go run ./cmd/meshgate
```

## Local development

```bash
go test ./...
go vet ./...
go test -race ./...
go run ./cmd/meshgate
```

With Docker:

```bash
docker build -t meshgate .
docker run --rm -p 8080:8080 \
  -e MESHGATE_ROUTE_PATH=/api \
  -e MESHGATE_ROUTE_TARGET=http://host.docker.internal:9000 \
  meshgate
```

## Testing strategy

The test suite is intended to exercise failure behavior, not just successful forwarding.

Current coverage includes:

- successful request forwarding
- upstream 5xx retry behavior
- retry safety for `POST`
- rate-limit rejection
- circuit opening and recovery
- race-detector coverage in CI

The benchmark harness is intentionally separate from functional tests. See [`benchmarks/README.md`](benchmarks/README.md).

## Performance evidence

No throughput or latency claim is published yet.

The project will publish benchmark results only after the request path and workload definitions stabilize. Reports will include commit SHA, Go version, hardware/OS, concurrency, request size, upstream behavior, throughput, p50/p95/p99 latency, error rate, retries, and breaker transitions.

This avoids presenting a single benchmark number without its workload and environment.

## Security posture

Security is treated as an engineering constraint, but this repository is **not** claiming production security readiness.

Current focus areas include request-size limits, bounded retries, upstream URL validation, non-root container execution, least-privilege CI permissions, and avoiding automatic retries of non-idempotent methods.

See [`SECURITY.md`](SECURITY.md) for the current disclosure and hardening scope.

## Deliberate limitations

The current build does **not** yet provide:

- distributed rate limiting
- service discovery
- weighted or adaptive load balancing
- Prometheus metrics
- OpenTelemetry tracing
- dynamic configuration reloads
- authentication/authorization middleware
- multi-route configuration files
- persistent configuration state
- production deployment manifests

These are explicit next steps, not implied capabilities.

## Roadmap

### Phase 1 — Reliability baseline

- [x] HTTP proxying
- [x] timeouts and deadlines
- [x] bounded retries with jitter
- [x] circuit breaker
- [x] request-size limits
- [x] local rate limiting
- [x] graceful shutdown
- [x] CI with race detection

### Phase 2 — Gateway engineering

- [ ] multiple route configuration
- [ ] health-aware load balancing
- [ ] upstream connection-pool controls
- [ ] deterministic routing tests
- [ ] structured request IDs

### Phase 3 — Observability and security

- [ ] Prometheus metrics
- [ ] OpenTelemetry traces
- [ ] authentication middleware
- [ ] upstream allow-list policy
- [ ] security-focused integration tests
- [ ] failure-injection scenarios

### Phase 4 — Evidence

- [ ] concurrency benchmark suite
- [ ] latency distributions under load
- [ ] controlled upstream failure experiments
- [ ] capacity and saturation analysis
- [ ] reproducible benchmark reports

## Engineering decisions

Important decisions are documented rather than hidden behind implementation details:

- [`docs/architecture.md`](docs/architecture.md)
- [`docs/reliability.md`](docs/reliability.md)
- [`benchmarks/README.md`](benchmarks/README.md)

The central rule is simple: **do not claim performance, scale, security, or reliability characteristics that have not been measured or tested.**

## License

MIT — see [`LICENSE`](LICENSE).
