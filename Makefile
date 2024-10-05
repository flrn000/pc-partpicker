.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N]' && read ans && [ $${ans:-N} = y ]

build: 
	@go build -o bin/partpicker cmd/main.go

.PHONY: test
test:
	@go test -v ./...

.PHONY: run
run: build
	@./bin/partpicker

.PHONY: psql
psql:
	psql ${DATABASE_URL}

.PHONY: migration
migration:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./cmd/migrate ${name}

.PHONY: up
up: confirm
	@echo 'Running up migrations'
	migrate -path ./cmd/migrate -database ${DATABASE_URL} up