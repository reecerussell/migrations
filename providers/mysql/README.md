# MySQL

Here is a basic example of setting up database migrations for MySQL To use the migrations tool with this provider, the connection string to the database must be defined in an environment variable `CONNECTION_STRING`.

### Connection String

The connection string must be in this format: `<username>:<password>@tcp(<host>)/<database>?parseTime=true`. This is defined by the MySQL driver, [github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql#dsn-data-source-name). It is worth noting, that your connection string will require the `parseTime` parameter to allow the process to parse MySQL times correctly.

### History

With the MySQL provider, migration history is stored in a table, named `__migration_history`. This table is used to record what migrations have been applied, and when.

The structure of the table is as follows:

| Column       | Type         | Allow Null | Auto Increment |
| ------------ | ------------ | ---------- | -------------  |
| id           | INT          | No         | YES            |
| name         | VARCHAR(255) | No         |                |
| date_applied | DATETIME     | No         |                |

### Configuration

When it comes to the migration configuration file, the `provider` property must be set to `mysql`.

Optionally, the migration history table name can be configured using the config map, `config`. A property named `historyTableName`, can be used to configure the history table name. If not set, the default value `__migration_history`, will be used.

```yaml
# migrations.yaml
provider: mysql
config:
    historyTableName: MyMigrations # default: __migration_history
migrations:
    - name: InitialCreation
      upFile: initialCreation.up.sql
      downFile: initialCreation.down.sql
```
