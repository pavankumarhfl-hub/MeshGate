# Contributing to MeshGate

Thanks for contributing.

## Engineering standard

MeshGate prioritizes correctness and evidence over feature count. Contributions should explain the engineering problem being solved and include tests for new behavior or failure modes.

### Before opening a change

```bash
go test ./...
go vet ./...
go test -race ./...
```

For behavior that affects performance or reliability, include a reproducible test or benchmark plan rather than an unsupported claim.

## Pull requests

A useful pull request should explain:

1. the problem;
2. the design and important trade-offs;
3. how the change was tested;
4. failure cases considered;
5. any remaining limitations.

Keep changes focused. Avoid introducing a dependency when the standard library is sufficient and the trade-off is not justified.

## Security

Do not disclose vulnerabilities in public issues. Follow [`SECURITY.md`](SECURITY.md) for responsible reporting.
