MAX_LINE_LENGTH := 100

.PHONY: web
web:
	cd textwire/example && go run main.go

.PHONY: shell
shell:
	go run repl/repl.go

.PHONY: test
test:
	clear
	go test ./...
	@echo "✅ All tests passed!"

.PHONY: fmt
fmt:
	@go fmt ./...
	echo "✅ Code formatted!"

.PHONY: line
line:
	@golines -w -m $(MAX_LINE_LENGTH) .
	echo "✅ Lines limited!"

.PHONY: lint
lint:
	@golangci-lint run
	echo "✅ Linting passed!"

.PHONY: check
check: test line fmt lint

.DEFAULT_GOAL := test
