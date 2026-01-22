CONTAINER := $(podman-compose run --rm)

.PHONY: test testp run runp shell shellp

test:
	echo "🚀 Running tests..."
	go test ./...
	echo "✅ Tests passed!"

testp:
	$(CONTAINER) make test

run:
	cd textwire/example && go run main.go

runp:
	$(CONTAINER) make run

shell:
	go run repl/repl.go

shellp:
	$(CONTAINER) make shell

.DEFAULT_GOAL := test
