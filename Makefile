SHELL := /usr/bin/env bash -o pipefail
.SHELLFLAGS := -ec

BUF_BREAKING_BRANCH ?= buf-breaking-base
BUF_BREAKING_TARGET ?= .git#branch=$(BUF_BREAKING_BRANCH)
LOCALBIN ?= $(shell pwd)/bin
GOBIN ?= $(LOCALBIN)
GOLANGCI_LINT = $(LOCALBIN)/golangci-lint

.PHONY: all
all: lint test

##@ Go targets

.PHONY: fmt
fmt: ## Run gofmt across the module
	go fmt ./...

.PHONY: lint
lint: golangci-lint ## Run golangci-lint
	"$(GOLANGCI_LINT)" run ./...

.PHONY: lint-fix
lint-fix: golangci-lint ## Run golangci-lint with autofix enabled
	"$(GOLANGCI_LINT)" run --fix ./...

.PHONY: lint-config
lint-config: golangci-lint ## Verify golangci-lint configuration
	"$(GOLANGCI_LINT)" config verify

.PHONY: test
test: ## Run Go unit tests
	go test $$(go list ./... | grep -v '^github.com/bubustack/tractatus/dist/')

.PHONY: tidy
tidy: ## Tidy go.mod / go.sum
	go mod tidy

##@ Protobuf targets

.PHONY: generate
generate: ## Regenerate tracked Go bindings (gen/go)
	@command -v buf >/dev/null 2>&1 || { echo "Error: buf CLI is not installed. See https://buf.build/docs/installation"; exit 1; }
	buf generate --template buf.gen.yaml

.PHONY: buf-lint
buf-lint: ## Run buf lint on proto definitions
	@command -v buf >/dev/null 2>&1 || { echo "Error: buf CLI is not installed. See https://buf.build/docs/installation"; exit 1; }
	buf lint

.PHONY: buf-breaking
buf-breaking: ## Ensure schema changes remain backwards compatible
	@command -v buf >/dev/null 2>&1 || { echo "Error: buf CLI is not installed. See https://buf.build/docs/installation"; exit 1; }
	@git rev-parse --verify refs/heads/$(BUF_BREAKING_BRANCH) >/dev/null 2>&1 || { echo "Error: local breaking baseline is unavailable. Run 'git fetch -f origin main:$(BUF_BREAKING_BRANCH)' after the first push to main."; exit 1; }
	buf breaking --against "$(BUF_BREAKING_TARGET)"

.PHONY: proto-dist
proto-dist: ## Build multi-language artifacts + descriptor bundle under dist/
	@command -v buf >/dev/null 2>&1 || { echo "Error: buf CLI is not installed. See https://buf.build/docs/installation"; exit 1; }
	rm -rf dist
	mkdir -p dist/release dist/release/descriptors
	buf generate --template buf.gen.yaml
	buf generate --template buf.gen.release.yaml
	buf build -o dist/release/descriptors/tractatus.bin
	rsync -a proto/ dist/release/proto/
	cp buf.yaml buf.gen.yaml buf.gen.release.yaml dist/release/
	cp LICENSE README.md dist/release/
	tar -czf dist/tractatus-protos.tar.gz -C dist/release .

##@ Tooling

GOLANGCI_LINT_VERSION ?= v2.11.4

.PHONY: golangci-lint
golangci-lint: $(GOLANGCI_LINT) ## Download golangci-lint locally if necessary.
$(GOLANGCI_LINT): $(LOCALBIN) .custom-gcl.yml
	$(call go-install-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/v2/cmd/golangci-lint,$(GOLANGCI_LINT_VERSION))
	@test -f .custom-gcl.yml && { \
		echo "Building custom golangci-lint with plugins..." && \
		$(GOLANGCI_LINT) custom --destination $(LOCALBIN) --name golangci-lint-custom && \
		mv -f $(LOCALBIN)/golangci-lint-custom $(GOLANGCI_LINT); \
	} || true

$(LOCALBIN):
	@mkdir -p $(LOCALBIN)

# go-install-tool will 'go install' any package with custom target and name of binary, if it doesn't exist
# $1 - target path with name of binary
# $2 - package url which can be installed
# $3 - specific version of package
define go-install-tool
@[ -f "$(1)-$(3)" ] && [ "$$(readlink -- "$(1)" 2>/dev/null)" = "$(1)-$(3)" ] || { \
set -e; \
package=$(2)@$(3) ;\
echo "Downloading $${package}" ;\
rm -f "$(1)" ;\
GOBIN="$(LOCALBIN)" go install $${package} ;\
mv "$(LOCALBIN)/$$(basename "$(1)")" "$(1)-$(3)" ;\
} ;\
ln -sf "$$(realpath "$(1)-$(3)")" "$(1)"
endef

##@ Utility

.PHONY: clean
clean: ## Remove build artifacts
	rm -rf dist

.PHONY: help
help: ## Print available targets
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_\-]+:.*##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0,5)}' $(MAKEFILE_LIST)
