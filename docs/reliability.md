# Reliability Model

MeshGate treats downstream failure as a normal operating condition.

## Timeout

Outbound requests use a bounded client timeout. A gateway should not wait indefinitely for a dependency.

## Retry

Only upstream failures represented as server errors or transport errors are retryable. Retries are bounded and use jittered backoff.

Retries must remain conservative because retrying can amplify load during an outage.

## Circuit breaker

Repeated failures move a route into an open state. While open, requests fail fast. After the cooldown, one request is allowed through as a half-open probe; success closes the breaker.

## Rate limiting

The current limiter provides a process-local per-client budget. It is not a distributed quota and should not be described as one.

## Failure scenarios to validate

- upstream unavailable
- upstream responds with repeated 5xx
- upstream becomes slow
- burst traffic exceeds the local limit
- concurrent requests contend on shared state
- gateway receives shutdown while requests are active
