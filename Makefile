.DEFAULT_GOAL := run

# Export all environments from within .env if it exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.PHONY: all

# List prints all targets in this makefile.
list:
	@make -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$'

# Cleans up the compiled binaries and tidies up the Go module dependencies.
clean:
	@rm -rf bin
	@go mod tidy

# Compiles the Go source code into a binary.
compile:
	@go build $(LDFLAGS) -o bin/emc main.go

# Builds the Go project and runs the resulting binary.
run:
	@go build -o bin/emc main.go
	@bin/emc $(ARGS)

# Runs all Go tests with the test_all tag, generates a coverage profile, and creates an HTML report of the coverage.
test:
	@mkdir -p coverage
	@go test ./... $(ARGS) --tags=test_all --coverprofile coverage/all.out -timeout 120s -parallel 1 -failfast -v
	@go tool cover -html=coverage/all.out -o coverage/all.html

# Runs unit tests with the unit tag, generates a coverage profile, and creates an HTML report of the coverage.
test.unit:
	@mkdir -p coverage
	@go test ./... $(ARGS) --tags=unit --coverprofile coverage/unit.out -timeout 120s -v
	@go tool cover -html=coverage/unit.out -o coverage/unit.html