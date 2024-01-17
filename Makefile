.PHONY: test
test:
	@echo "🚀 Running tests..."
	@go test -v ./...
	@echo "✅ Tests passed!"

.PHONE: run
run:
	@cd example && go run main.go

.DEFAULT_GOAL := test