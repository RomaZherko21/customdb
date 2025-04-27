run:
	go run ./cmd/main.go -mode console

run-http:
	go run ./cmd/main.go -mode http -port 8080

test:
	go test ./internal/... -v

e2e:
	go test ./e2e/... -v

generate:
	swag init -g ./cmd/main.go -o ./docs
