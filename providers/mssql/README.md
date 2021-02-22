# SQL Server

Here is a basic example of setting up database migrations for SQL Server. To use the migrations tool with this provider, the connection string to the database must be defined in an environment variable `CONNECTION_STRING`.

### Connection String

The connection string must be in this format: `sqlserver://<username>:<password>@<host>?database=<database>`. This is defined by the SQL Server driver, [github.com/denisenkom/go-mssqldb](https://github.com/denisenkom/go-mssqldb#the-connection-string-can-be-specified-in-one-of-three-formats).

### History

With the SQL Server provider, migration history is stored in a SQL table, named `__MigrationHistory`. This table is used to record what migrations have been applied, and when.

The structure of the table is as follows:

| Column      | Type         | Allow Null |               |
| ----------- | ------------ | ---------- | ------------- |
| Id          | INT          | No         | IDENTITY(1,1) |
| Name        | VARCHAR(255) | No         |               |
| DateApplied | DATETIME     | No         |               |

<style>
table {
    width: 100%
}
</style>

### Configuration

When it comes to the migration configuration file, the `provider` property must be set to `mssql`.

Optionally, the migration history table name can be configured using the config map, `config`. A property named `historyTableName`, can be used to configure the history table name. If not set, the default value `__MigrationHistory`, will be used.

```yaml
# migrations.yaml
provider: mssql
config:
    historyTableName: MyMigrations # default: __MigrationHistory
migrations:
    - name: InitialCreation
      upFile: initialCreation.up.sql
      downFile: initialCreation.down.sql
```
