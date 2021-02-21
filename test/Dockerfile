FROM golang
WORKDIR /go/src/github.com/reecerussell/migrations
COPY . .

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

RUN go mod download

RUN chmod +x test/entrypoint.sh
ENTRYPOINT [ "./test/entrypoint.sh" ]