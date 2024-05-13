build:
	@go build -o bin/ecom cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/ecom 

install:
	@go mod tidy	 