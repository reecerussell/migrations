[![Go Report Card](https://goreportcard.com/badge/github.com/reecerussell/migrations)](https://goreportcard.com/report/github.com/reecerussell/migrations)
[![codecov](https://codecov.io/gh/reecerussell/migrations/branch/master/graph/badge.svg?token=MRH242FDJE)](https://codecov.io/gh/reecerussell/migrations)
![Actions](https://github.com/reecerussell/migrations/actions/workflows/release.yaml/badge.svg)
[![Go Docs](https://godoc.org/github.com/reecerussell/migrations?status.svg)](https://godoc.org/github.com/reecerussell/migrations)

# Migrations

Migrations is a simple tool used to manage database migrations. With support to apply and rollback changes, Migrations is designed to run in both CI/CD pipelines and from a dev machine.

## Using Docker

```bash
# ./
docker run --rm --name migrations \
    --mount type=bind,source="$(pwd)/example",target=/migrations \
    -e CONNECTION_STRING="<connection string>" \
    migrations up --context /migrations
```