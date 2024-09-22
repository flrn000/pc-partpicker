build: 
	@go build -o bin/partpicker cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/partpicker