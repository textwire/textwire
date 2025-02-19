.PHONY: test
test:
	echo "🚀 Running tests..."
	go test ./...
	echo "✅ Tests passed!"

.PHONE: run
run:
	clear
	@cd textwire/example && go run main.go

.PHONE: shell
shell:
	go run repl/repl.go

.PHONE: push
push: test
	echo "🚀 Pushing to GitHub..."
	git push
	echo "✅ Pushed to GitHub!"

.DEFAULT_GOAL := test
