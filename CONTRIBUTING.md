# Contributing to Tractatus

Tractatus houses the transport proto definitions, generated Go bindings, and helper packages that every BubuStack component relies on. Thanks for helping us keep the contract coherent, current, and multi-language friendly.

## Reporting bugs

1. Search [existing issues](https://github.com/bubustack/tractatus/issues?q=is%3Aissue) before filing a new one.
2. Include the following details when filing a bug:
   - Which generated artifact or helper package (`envelope`, `transport`, etc.) is affected.
   - The versions of Tractatus, bobrapet, and any SDK/Engram using it.
   - Reproduction steps (proto snippet, Go sample, or CLI commands).
   - Whether the issue impacts wire compatibility with existing deployments.
3. Tag the issue with `kind/bug`, an appropriate `area/*` label (`area/proto`, `area/envelope`, `area/generated`, `area/release`, or `area/tooling`), and a priority label if known.

## Requesting enhancements

- Use the [feature request template](https://github.com/bubustack/tractatus/issues/new?template=feature_request.md) to propose new fields, events, or generated artifacts.
- Explain the workflow or SDK that needs the change, including backward-compatibility or rollout considerations.
- For schema changes, describe the Buf breaking policy impact and whether new releases of downstream SDKs/operators are required.

## Pull requests

1. Fork the repo, branch from `main`, and keep each PR focused.
2. Run the quality gates before requesting review:
   ```bash
   make lint          # golangci-lint for Go helpers
   make test          # Go unit tests
   make buf-lint      # Proto style checks via buf
   make generate      # Regenerate Go bindings if proto/schema changed
   ```
3. When changing protos, fetch the compatibility baseline and run `make buf-breaking` locally when you need to measure drift against `main`:
   ```bash
   git fetch origin main:main
   make buf-breaking
   ```
4. Update `README.md` and release notes when you add new event types or streams.
5. Fill out the PR template with commands you ran, test output, and linked issues (e.g., `Fixes #123`).

## Development workflow

### Prerequisites

- Go 1.25.1 or newer
- [Buf CLI](https://buf.build/docs/installation) for linting, breaking-change checks, and code generation
- `make`

### Useful targets

- `make lint` / `make lint-fix`
- `make test`
- `make buf-lint` / `make buf-breaking`
- `make generate` — regenerates tracked Go bindings under `gen/go`
- `make proto-dist` — builds multi-language artifacts and descriptor bundles under `dist/`

### Protobuf-specific rules

- Treat `proto/transport/v1/*` as the current shared API surface.
- Before the first public release, latest `main` is the source of truth. Coordinated breaking changes are allowed if every downstream consumer is updated in the same rollout.
- After the first public release, prefer additive evolution and avoid breaking wire changes in `transport.v1`.
- Never edit files under `gen/go/` by hand. Regenerate them with `make generate`.
- When removing or renaming fields in a future breaking release, reserve field numbers and names explicitly.

## Commit message conventions

We follow [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) so release-please can build release notes and version bumps consistently.

Examples:

- `feat: add connector codec selection action`
- `fix: correct envelope MIME validation`
- `docs: update README hub quickstart`
- `chore: refresh buf plugin versions`

## Code of Conduct

Participation in this project is governed by the [Contributor Covenant Code of Conduct](./CODE_OF_CONDUCT.md). Report unacceptable behavior to [conduct@bubustack.com](mailto:conduct@bubustack.com) or via the Discussions moderation channel.
