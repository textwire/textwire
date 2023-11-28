.PHONY: test
test:
	@echo "🚀 Running tests..."
	@go test -v ./...
	@echo "✅ Tests passed!"

.DEFAULT_GOAL := test