build:
	@go build -o bin/api cmd/main/main.go


run: build
	@./bin/api

test:
	@go test -v ./...

seed:
	@go run cmd/seed/seed.go

workers:
	@go build -o bin/worker/message cmd/worker/message_worker.go
	@./bin/worker/message

message-test:
	@go run pkg/workers/message/cmd/test/main.go