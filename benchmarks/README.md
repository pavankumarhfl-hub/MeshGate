# Benchmarking

Benchmark results are intentionally kept out of the repository until the request path stabilizes.

A valid report must record:

- commit SHA
- Go version
- operating system and CPU
- workload and request size
- concurrency
- duration and warm-up
- upstream behavior
- throughput
- p50/p95/p99 latency
- error rate
- retry count
- circuit-breaker transitions

Do not publish a single throughput number without the workload and environment that produced it.
