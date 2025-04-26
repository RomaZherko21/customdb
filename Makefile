run:
	go run ./cmd/main.go -mode console

run-http:
	go run ./cmd/main.go -mode http -port 8080

test:
	go test ./... -v

swag-gen:
	swag init -g ./cmd/main.go -o ./docs
