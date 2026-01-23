.PHONY: web test shell fmt lint check test-count

web:
	cd textwire/example && go run main.go

shell:
	go run repl/repl.go

test:
	echo "🚀 Running tests..."
	go test ./...
	@echo "✅ $$(make -s test-count) tests pass"

fmt:
	echo "🔧 Formatting code..."
	go fmt ./...
	echo "✅ Code formatted!"

lint:
	echo "🔍 Running linter..."
	golangci-lint run
	echo "✅ Linting passed!"

check: fmt lint test

test-count:
	@go test -json ./... | jq -s '[.[] | select(.Action == "run" and .Test != null)] | length'

.DEFAULT_GOAL := test
