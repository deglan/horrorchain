APP_NAME := blockhorror

.PHONY: run build fmt lint clean assets

run:
	go run .

tidy:
	go mod tidy

build:
	go build -o $(APP_NAME) main.go

fmt:
	go fmt ./...

lint:
	golangci-lint run

clean:
	rm -f $(APP_NAME)

