MAX_LINE_LENGTH := 100

.PHONY: build
build:
	@clear || true
	cd textwire/example && go build main.go

.PHONY: repl
repl:
	@clear || true
	@go run repl/repl.go

.PHONY: cover
cover:
	go test -coverprofile=coverage.html
	go tool cover -html=coverage.html
	rm coverage.html

.PHONY: test
test:
	@clear || true
	go test ./...
	@echo "✅ All tests passed!"

.PHONY: fmt
fmt:
	@go fmt ./...
	@echo "✅ Code formatted!"

.PHONY: line
line:
	@golines -w -m $(MAX_LINE_LENGTH) .
	@echo "✅ Lines limited!"

.PHONY: lint
lint:
	@golangci-lint run
	@echo "✅ Linting passed!"

.PHONY: bench
bench:
	@go test ./textwire/example -bench=BenchmarkTestProject -benchmem

.PHONY: todo
todo:
	@if grep -I --exclude="Makefile" --exclude-dir=".git" -r TODO .; then \
		echo "❌ Found TODOs" >&2; \
		exit 1; \
	fi

.PHONY: check
check: test line fmt lint todo

.DEFAULT_GOAL := test
