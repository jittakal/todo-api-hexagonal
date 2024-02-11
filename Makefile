.PHONY: all clean generate run-http run-grpc test coverage

# Define variables
GOPATH := $(shell go env GOPATH)
SWAG_CMD := $(GOPATH)/bin/swag
PROTOC_CMD := protoc
TEST_FLAGS := -race -v -covermode=atomic -coverprofile=coverage.out

# Default target
all: generate run-http

clean:
	# rm -rf ./docs

# Generate Swagger documentation
generate:
	$(SWAG_CMD) init -g ./internal/adapters/handler/rest/http_task_handler.go

# Run the HTTP and gRPC server
run:
	go run cmd/main.go

# Run HTTP server
run-http:
	go run cmd/main.go --handler http 

# Run gRPC server
run-grpc:
	go run cmd/main.go --handler grpc

# Generate Go files from Protobuf
generate-proto:
	$(PROTOC_CMD) --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/proto/task.proto

# Run tests
test:
	go test $(TEST_FLAGS) ./...

# Generate code coverage report
coverage:
	go test $(TEST_FLAGS) ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out
