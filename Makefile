all: deps build
test: run-tests

deps:
	go mod download
	go mod verify

generate:
	go generate mock/mock.go

run-tests:
	docker-compose up --build --exit-code-from tests

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o migrations cmd/main.go

build-app:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /app/migrations cmd/main.go