FROM golang
WORKDIR /go/src/github.com/reecerussell/migrations
COPY . .

RUN go mod download

ENV CONNECTION_STRING=sqlserver://sa:Password_123@mssql?database=test

RUN chmod +x test/entrypoint.sh
ENTRYPOINT [ "./test/entrypoint.sh" ]