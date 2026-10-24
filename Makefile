MAX_LINE_LENGTH := 100

OCI ?= podman-compose
RUN = $(OCI) run --rm -T checks
RUN_IT= $(OCI) run --rm checks
IMAGE ?= registry.serhiicho.com/textwire:checks

.PHONY: image
image:
	@$(OCI) build

.PHONY: repl
repl: image
	@clear || true
	@$(RUN_IT) go run pkg/repl/repl.go

.PHONY: cover
cover: image
	@$(RUN) go test -coverprofile=coverage.out
	@$(RUN) go tool cover -html=coverage.out -o coverage.html
	@rm -f coverage.out
	@echo "✅ Coverage written to coverage.html"

.PHONY: test
test: image
	@clear || true
	@$(RUN) go test -count=1 ./...
	@echo "✅ All tests passed!"

.PHONY: fmt
fmt: image
	@$(RUN) go fmt ./...
	@echo "✅ Code formatted!"

.PHONY: line
line: image
	@$(RUN) golines -w -m $(MAX_LINE_LENGTH) .
	@echo "✅ Lines limited!"

.PHONY: lint
lint: image
	@$(RUN) golangci-lint run
	@echo "✅ Linting passed!"

.PHONY: todo
todo: image
	@if $(RUN) grep -I --exclude="Makefile" --exclude-dir=".git" -r TODO .; then \
		echo "❌ Found TODOs" >&2; \
		exit 1; \
	fi

.PHONY: check
check: test line fmt lint todo

.DEFAULT_GOAL := test
