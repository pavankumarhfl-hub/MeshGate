# Architecture

MeshGate is organized as a small gateway runtime with explicit reliability boundaries.

```text
Client
  |
  v
HTTP Server
  |
  +--> rate limiter
  |
  +--> route lookup
  |
  v
Proxy
  |
  +--> timeout
  +--> retry + jitter
  +--> circuit breaker
  |
  v
Upstream service
```

## Current boundaries

- `internal/gateway`: request routing, health endpoints and middleware composition.
- `internal/proxy`: outbound HTTP requests, timeout handling and retry behavior.
- `internal/ratelimit`: synchronized per-client request limiting.
- `internal/breaker`: closed/open/half-open circuit state.
- `internal/config`: runtime defaults.

## Design principles

1. Keep the first implementation dependency-light so behavior is easy to inspect.
2. Keep reliability mechanisms isolated so they can be tested independently.
3. Treat timeouts and bounded retries as correctness controls, not performance features.
4. Measure behavior before making performance claims.

## Known limitations

The current implementation is intentionally an early systems-engineering baseline. It does not yet provide dynamic configuration, service discovery, weighted load balancing, distributed rate limiting, persistent state, metrics/tracing exporters, or a complete graceful-shutdown path. Those are tracked as subsequent engineering work rather than claimed capabilities.
