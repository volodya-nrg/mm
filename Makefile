PROJECT_NAME=mm
GIT_COMMIT=$(shell git rev-parse --short HEAD)
VERSION_FLAG="-X main.version=commit:$(GIT_COMMIT)"

.PHONY: echo_version
echo_version:
	@echo commit:$(GIT_COMMIT)

.PHONY: build
build:
	go build -ldflags=$(VERSION_FLAG) -o ./bin/$(PROJECT_NAME) ./cmd/server

.PHONY: run_server
run_server: build
	./bin/$(PROJECT_NAME)

.PHONY: transport_tests
transport_tests:
	go test -v ./tests/transport

.PHONY: lint
lint:
	cd scripts && bash ./lint_docker.sh
