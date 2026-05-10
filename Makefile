.PHONY: build run-api run-tcp run-udp run-grpc run-all proto init seed test clean fmt tidy

## Build all binaries
build:
	go build ./cmd/...

## Run individual servers
run-api:
	go run ./cmd/api-server

run-tcp:
	go run ./cmd/tcp-server

run-udp:
	go run ./cmd/udp-server

run-grpc:
	go run ./cmd/grpc-server

## Run all servers in background (dev convenience)
run-all:
	@echo "Starting all MangaHub servers..."
	go run ./cmd/api-server  &
	go run ./cmd/tcp-server  &
	go run ./cmd/udp-server  &
	go run ./cmd/grpc-server &
	@echo "All servers started. Use 'make clean' to stop."

## Generate protobuf Go code
## Requires: protoc + protoc-gen-go + protoc-gen-go-grpc
proto:
	protoc \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/manga.proto

## Initialize ~/.mangahub/ config and directories (run once)
init:
	go run ./cmd/init

## Seed manga database from data/manga.json
seed:
	go run ./cmd/seed

## Run all tests
test:
	go test ./... -v -race

## Format code
fmt:
	gofmt -w .

## Sync go.sum
tidy:
	go mod tidy

## Kill all mangahub server processes
clean:
	-pkill -f "api-server"  2>/dev/null || true
	-pkill -f "tcp-server"  2>/dev/null || true
	-pkill -f "udp-server"  2>/dev/null || true
	-pkill -f "grpc-server" 2>/dev/null || true
	@echo "All servers stopped."
