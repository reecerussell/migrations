FROM golang:alpine AS base

RUN apk update && apk add --no-cache make

ENV USER=app
ENV UID=10001

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistant" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

FROM base AS build
WORKDIR /go/src/github.com/reecerussell/migrations
COPY . .

RUN make

FROM scratch
WORKDIR /app

COPY --from=base /etc/passwd /etc/passwd
COPY --from=base /etc/group /etc/group
COPY --from=build /app/migrations migrations

RUN mkdir -p /migrations

USER ${UID}

ENTRYPOINT [ "./migrations" ]
CMD [ "up", "-context", "/migrations" ]