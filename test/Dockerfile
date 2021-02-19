FROM golang

RUN apt-get update && apt-get make netcat

COPY . /go/src/github.com/reecerussell/migrations

RUN make install-deps

ENV CONNECTION_STRING=sqlserver://sa:Password_123@mssql?database=test

RUN sleep 30
ENTRYPOINT [ "make", "run-tests" ]