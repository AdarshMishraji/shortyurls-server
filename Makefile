dev:
	gomon ./cmd/main.go

start:
	go ./cmd/main.go

build:
	go build -o ./bin/main ./cmd/main.go

build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o ./bin/main ./cmd/main.go

