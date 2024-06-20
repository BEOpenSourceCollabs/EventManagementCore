.DEFAULT_GOAL := run

# Export all environments from within .env if it exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# Define vars
DB_DSN="postgres://$(DATABASE_USER):$(DATABASE_PASSWORD)@$(DATABASE_HOST):$(DATABASE_PORT)/$(DATABASE_NAME)?sslmode=$(DATABASE_SSL_MODE)"

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


# Print current version of database
migrate.version:
	@migrate -path=./migrations -database=$(DB_DSN) version

# Create a new migration
# Usage: `make migrate.create NAME=<your-migration-name>`
migrate.create:
	@migrate create -seq -ext=.sql -dir=./migrations ${NAME}

# Apply all migrations from first to latest
migrate.all.up:
	@migrate -path=./migrations -database=$(DB_DSN) up

# Remove all migrations
migrate.all.down:
	@migrate -path=./migrations -database=$(DB_DSN) down

# Go to a specific version of database 
# Usage `make migrate.goto V=4`
# If current version is 2 then #3 and #4 `Up` will be applied on the other hand if current version
# is 6 then #6 and #5 will be `Down` will be applied. 
migrate.goto:
	@migrate -path=./migrations -database=$(DB_DSN) goto $(V)

# Force a specific version of database, useful to force database to specific version 
# after manually fixing migration issues. (only use this command if there were syntax errors
# in your migration files and you are not sure if migration was applied fully or partially, in that case
# you need to manually fix issues and then manually set the version upto which you have fixed issues.)
# Usage `make migrate.force V=1` sets the current version of database to 1
migrate.force:
	@migrate -path=./migrations -database=$(DB_DSN) force $(V)