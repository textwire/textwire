.PHONY: test
test:
	@echo "🚀 Running tests..."
	@go test -v ./...
	@echo "✅ Tests passed!"

.PHONE: run
run:
	@cd example && go run main.go

.PHONE: push
push:
	@echo "🚀 Pushing to GitHub..."
	make test
	git push
	@echo "✅ Pushed to GitHub!"

.DEFAULT_GOAL := test