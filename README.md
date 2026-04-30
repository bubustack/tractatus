# 📦 Tractatus — BubuStack Transport Contracts

[![Go Reference](https://pkg.go.dev/badge/github.com/bubustack/tractatus.svg)](https://pkg.go.dev/github.com/bubustack/tractatus)
[![Go Report Card](https://goreportcard.com/badge/github.com/bubustack/tractatus)](https://goreportcard.com/report/github.com/bubustack/tractatus)

Tractatus contains the canonical transport proto files, generated Go bindings, and helper packages (`envelope`, `transport`) that every bobrapet operator, SDK, Engram, and connector relies on. Keeping these contracts in a dedicated repo makes it easy to lint, version, and publish checked-in Go bindings plus release bundles for other language targets.

## Repository layout

| Path | Purpose |
| --- | --- |
| `proto/transport/v1/` | Protobuf definitions for hub RPCs, binary frames, and transport envelopes. |
| `gen/go/proto/transport/v1/` | Checked-in Go bindings generated via `buf generate` (kept in sync with `go.mod`). |
| `envelope/` | Helpers for marshaling/unmarshaling the JSON envelope format used inside BinaryFrame payloads. |
| `transport/` | Shared constants for stream type identifiers (`speech.audio.v1`, `speech.transcript.delta`, etc.). |
| `validation/` | Runtime helper for applying Protovalidate rules embedded in generated transport bindings. |

## Current services

- `HubService.Process` is the hub-facing bidirectional stream for `DataPacket` traffic.
- `TransportConnectorService.Data`, `Control`, and `HubPush` are the connector-side bidirectional streams used by SDKs and sidecars.

## Quick start (Go)

```bash
go get github.com/bubustack/tractatus
```

```go
import (
    "encoding/json"

    transportpb "github.com/bubustack/tractatus/gen/go/proto/transport/v1"
    "github.com/bubustack/tractatus/envelope"
)

// conn is a grpc.ClientConnInterface and ctx is already in scope.
client := transportpb.NewHubServiceClient(conn)
stream, _ := client.Process(ctx)

env := &envelope.Envelope{
    Payload: json.RawMessage(`{"message":"hello"}`),
}
frame, _ := envelope.ToBinaryFrame(env)

_ = stream.Send(&transportpb.ProcessRequest{
    Packet: &transportpb.DataPacket{
        Frame: &transportpb.DataPacket_Binary{Binary: frame},
    },
})
```

The Go bindings are checked in to this repo, so consumers **do not** need protoc or buf installed.

### Runtime validation

The proto bindings include Protovalidate annotations for transport payload
sizes, metadata maps, repeated fields, media bounds, and chunk envelope
invariants. Consumers can validate generated messages through the project helper:

```go
import "github.com/bubustack/tractatus/validation"

if err := validation.Validate(msg); err != nil {
    // reject or report the invalid transport message
}
```

Validation is explicit. gRPC streams do not reject messages automatically unless
the caller invokes this helper or a Protovalidate runtime at the boundary.

### Authentication, TLS, and stream bounds

Tractatus defines message shape; it does not authenticate callers by itself.
Production hubs and connectors should run these streams over authenticated
transport channels and bind the presented identity to claimed metadata such as
namespace, StoryRun, step, or connector identity. mTLS SAN checks, workload
identity, or an equivalent deployment-specific mechanism should happen before
accepting packets.

Implementations should enforce runtime limits outside the protobuf schema:

- gRPC maximum message size aligned with the payload bounds in the proto.
- Keepalive policy for idle but healthy streams.
- Rate limiting or coalescing for high-frequency control directives,
  acknowledgements, heartbeats, and telemetry.
- Maximum stream or session duration where required by the deployment.
- Rejection of traffic from stale or superseded connectors after topology
  cutover.

### Envelope safety contract

The `envelope` helpers enforce a bounded JSON envelope contract for
`BinaryFrame` payloads:

- `envelope.Unmarshal` rejects payloads larger than `envelope.MaxEnvelopeSize`
  (10 MiB).
- `envelope.Marshal` rejects invalid JSON in `Payload` and `Inputs`.
- `TimestampMs` must be non-negative, and `FromBinaryFrame` rejects frame
  timestamps that cannot fit in the Go envelope timestamp field.

## Development workflow

```bash
make lint          # golangci-lint for helper packages
make test          # Go unit tests
make buf-lint      # Buf lint on proto definitions
make buf-breaking  # Compatibility diff against the local main baseline
make generate      # Regenerate tracked Go bindings under gen/go
```

All commands rely on the Buf CLI and Go 1.26+. The lint (`.github/workflows/lint.yml`) and test (`.github/workflows/test.yml`) workflows run the same checks on every push and pull request.

## Publishing multi-language artifacts

When a proto change lands, run:

```bash
make proto-dist
```

This target:

1. Regenerates the Go bindings under `gen/go`.
2. Uses `buf.gen.release.yaml` to emit ecosystem-ready packages:
   - `dist/release/go/` — Go module bindings (proto + gRPC).
   - `dist/release/ts/` — TypeScript outputs via ts-proto (gRPC JS clients).
   - `dist/release/python/` — Python protobuf + gRPC stubs.
3. Builds a descriptor set (`dist/release/descriptors/tractatus.bin`).
4. Copies the raw `proto/` tree, Buf templates, README, and LICENSE into `dist/release/`.
5. Bundles everything into `dist/tractatus-protos.tar.gz` for manual distribution when we build non-Go SDKs.

Releases currently focus on the Go module (indexed automatically by pkg.go.dev). When a GitHub Release is published, the `Publish Proto Bundle` workflow builds `dist/tractatus-protos.tar.gz` and uploads it to the release automatically. Extending `buf.gen.release.yaml` with new plugins is all that’s required to add more language targets.

## Releasing

We use [release-please](https://github.com/googleapis/release-please) to manage version bumps and GitHub Releases. On pushes to `main`, the workflow updates or opens the release PR; merging that PR cuts the tag and GitHub Release. `pkg.go.dev` picks up the new tag automatically, and the proto bundle is uploaded by the release workflow chain.

## Community & support

- [Contributing guide](./CONTRIBUTING.md)
- [Support](./SUPPORT.md)
- [Discord](https://discord.gg/dysrB7D8H6)
- [Security policy](./SECURITY.md)
- [Code of Conduct](./CODE_OF_CONDUCT.md)

## 📄 License

Copyright 2025 BubuStack.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
