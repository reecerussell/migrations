all: build run
test: generate run-tests

install-deps:
	go mod download

generate:
	go generate mock/mock.go

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o migrations cmd/main.go

run-tests:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go test ./... -cover -coverprofile=/tests/coverage.out

run:
	./migrations