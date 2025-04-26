run:
	go run ./cmd/main.go -mode console

run-http:
	go run ./cmd/main.go -mode http -port 80

test:
	go test ./... -v

