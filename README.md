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