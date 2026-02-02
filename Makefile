MAX_LINE_LENGTH := 100

.PHONY: dev
dev:
	clear
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

.PHONY: todo
todo:
	@if grep --exclude="Makefile" -r TODO .; then \
		echo "❌ Found TODOs" >&2; \
		exit 1; \
	fi

.PHONY: check
check: test line fmt lint todo

.DEFAULT_GOAL := test
