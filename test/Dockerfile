FROM golang
WORKDIR /go/src/github.com/reecerussell/migrations
COPY . .

RUN go mod download
RUN go mod verify

RUN chmod +x test/entrypoint.sh
ENTRYPOINT [ "./test/entrypoint.sh" ]