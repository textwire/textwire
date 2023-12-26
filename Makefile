.PHONY: test
test:
	@echo "🚀 Running tests..."
	@go test -v ./...
	@echo "✅ Tests passed!"

.PHONY: cli
cli:
	@echo "🚀 Running CLI..."
	@go run cli/cli.go

.DEFAULT_GOAL := test