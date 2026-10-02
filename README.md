# MeshGate

**A reliability-focused HTTP API gateway written in Go.**

MeshGate is an engineering project for studying how gateways behave when downstream services are slow, unavailable, overloaded, or returning errors. The project emphasizes explicit failure handling, bounded retries, circuit breaking, concurrency safety, testing, and measurable behavior.

> MeshGate is an engineering project, not a claim of production readiness. Capabilities are documented according to what is implemented and tested.

## Why this project exists

A reverse proxy that forwards HTTP requests is easy to build. A useful gateway must also define what happens when dependencies fail.

MeshGate is being developed around those questions:

- How long should a request wait?
- Which failures are safe to retry?
- How do retries avoid amplifying an outage?
- When should a dependency be temporarily isolated?
- How should concurrent clients share gateway resources?
- How do we test failure paths rather than only happy paths?
- Which performance claims can actually be supported by measurements?

## Current capabilities

- HTTP routing to configured upstreams
- Bounded outbound request timeouts
- Bounded retries for transport and upstream 5xx failures
- Jittered retry backoff
- Process-local per-client rate limiting
- Circuit breaker with closed/open/half-open states
- Health and readiness endpoints
- JSON route listing endpoint
- Unit and integration tests
- Go race-detector CI
- Minimal non-root production container

## Architecture

```text
                    +-------------------+
Client ------------>|     MeshGate      |
                    |                   |
                    | Rate Limiter      |
                    | Route Lookup      |
                    |       |           |
                    |       v           |
                    |   HTTP Proxy      |
                    |   /  |  \         |
                    | timeout retries   |
                    |        |          |
                    |  circuit breaker  |
                    +--------|----------+
                             |
                             v
                       Upstream service
```

See [`docs/architecture.md`](docs/architecture.md) for the current component boundaries and [`docs/reliability.md`](docs/reliability.md) for failure semantics.

## Run locally

```bash
go run ./cmd/meshgate
```

The default listener is `:8080`.

```bash
curl http://localhost:8080/healthz
curl http://localhost:8080/readyz
curl http://localhost:8080/routes
```

The current baseline registers `/example` against `http://127.0.0.1:9000`. A later configuration layer will replace this development-only default with explicit route configuration.

## Test

```bash
go test ./...
go test -race ./...
go vet ./...
```

CI runs tests, `go vet`, and the race detector across supported Go versions.

## Engineering decisions

### Go standard library first

The initial implementation deliberately minimizes dependencies. HTTP behavior, synchronization and failure handling remain visible in the source and easy to test.

### Bounded retries

Retries are limited because an unhealthy upstream can otherwise turn a small failure into a larger traffic storm. Jitter reduces synchronized retry timing.

### Circuit breaking

Repeated failures should eventually fail fast instead of continuously consuming gateway and upstream resources.

### Process-local rate limiting

The first implementation keeps state in memory to make semantics deterministic and dependency-free. This is intentionally **not** a distributed rate limiter.

## Performance

No throughput or latency numbers are published yet. Benchmarking will be added only after the request path and failure semantics stabilize. Results will include workload, concurrency, hardware/runtime information, and methodology so they can be reproduced.

## Roadmap

### Phase 1 — reliability baseline

- [x] HTTP gateway runtime
- [x] bounded timeouts
- [x] retries with jitter
- [x] circuit breaker
- [x] process-local rate limiting
- [x] health/readiness endpoints
- [x] unit/integration tests
- [x] race-detector CI
- [x] container image

### Phase 2 — gateway engineering

- [ ] explicit YAML/JSON configuration
- [ ] multiple upstreams per route
- [ ] round-robin and least-connections load balancing
- [ ] active/passive health checks
- [ ] graceful shutdown and connection draining
- [ ] request/response size limits
- [ ] structured access logging

### Phase 3 — observability and security

- [ ] Prometheus-compatible metrics
- [ ] OpenTelemetry tracing
- [ ] request IDs and correlation
- [ ] authentication middleware interface
- [ ] security headers and policy documentation
- [ ] audit events

### Phase 4 — evidence

- [ ] concurrency benchmarks
- [ ] latency distributions
- [ ] retry/circuit-breaker failure experiments
- [ ] load-shedding experiments
- [ ] documented capacity limits
- [ ] reproducible benchmark reports

## Security

Security is treated as an engineering constraint. The project will document trust boundaries, authentication assumptions, header handling, resource limits, dependency policy and abuse cases before claiming production readiness.

Please see the repository's security policy when it is added.

## Project status

**Early engineering build.** The repository is intentionally incomplete. Missing production features are tracked explicitly instead of being presented as implemented capabilities.

## License

MIT
