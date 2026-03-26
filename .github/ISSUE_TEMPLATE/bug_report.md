---
name: "Bug report"
about: "Report a reproducible issue in the tractatus contracts, generated bindings, or helper packages"
labels: ["kind/bug"]
---

## Area
- [ ] Proto schema (`proto/transport/v1/*`)
- [ ] Generated Go bindings (`gen/go/...`)
- [ ] Envelope helper package
- [ ] Transport stream constants
- [ ] Release / CI / Buf tooling
- [ ] Documentation

## What happened?
Tell us what broke, what you expected, and what you observed instead.

## Minimal reproduction
1. Proto/helper API snippet or generated type involved
2. Commands you ran (`make`, `go test`, `buf lint`, etc.)
3. Consumer(s) affected (`bobrapet`, `bobravoz-grpc`, `bubu-sdk-go`, specific engram, etc.)

```
make buf-lint
make generate
go test ./...
```

## Compatibility impact
- Does this change break wire compatibility?
- Are existing generated consumers affected?
- If yes, which repos/packages fail?

## Additional context
Anything else we should know? Link failing PRs, downstream repos, or relevant release notes.
