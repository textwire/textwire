.PHONY: web
web:
	cd textwire/example && go run main.go

.PHONY: shell
shell:
	go run repl/repl.go

.PHONY: test
test:
	echo "🚀 Running tests..."
	go test ./...
	@echo "✅ All tests passed!"

.PHONY: fmt
fmt:
	echo "🔧 Formatting code..."
	go fmt ./...
	echo "✅ Code formatted!"

.PHONY: lint
lint:
	echo "🔍 Running linter..."
	golangci-lint run
	echo "✅ Linting passed!"

.PHONY: check
check: fmt lint test
.DEFAULT_GOAL := test
