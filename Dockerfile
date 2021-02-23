FROM golang:alpine AS build

RUN apk update && apk add make

WORKDIR /go/src/github.com/reecerussell/migrations
COPY . .

RUN make deps

# Run unit tests except for provider specific tests, so those
# will be covered in integration testing stages.
RUN go test . ./providers

RUN make build-app

FROM scratch
WORKDIR /app

COPY --from=build /app/migrations .

ENTRYPOINT [ "./migrations" ]