build:
	go build -o bin/api cmd/main/main.go

run: build
	./bin/api

test:
	go test -v ./...

