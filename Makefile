run:
	go run ./cmd/main.go -mode console

run-http:
	go run ./cmd/main.go -mode http -port 8080


test:
	go test ./internal/... -v

test-e2e:
	go test ./e2e/... -v

generate:
	swag init -g ./cmd/main.go -o ./docs


.PHONY: protoGen
protoGen: # Generating client and server code (.pb which contains all the protocol buffer code to populate, serialize, and retrieve request and response message types; _grpc.pb An interface type for clients and servers)
	protoc --go_out=./internal/disk_manager/generated \
    --go-grpc_out=./internal/disk_manager/generated \
    ./internal/disk_manager/Page.proto