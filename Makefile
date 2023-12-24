.PHONY: test
test:
	@echo "🚀 Running tests..."
	@go test -v ./...
	@echo "✅ Tests passed!"

.PHONY: run
run:
	@echo "🚀 Running app..."
	@go run main.go

.DEFAULT_GOAL := test