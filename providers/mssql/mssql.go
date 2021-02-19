package mssql

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/reecerussell/migrations"
	"github.com/reecerussell/migrations/providers"

	// MSSQL driver
	_ "github.com/denisenkom/go-mssqldb"
)

const historyTableName = "__MigrationHistory"

var connectionString = os.Getenv("CONNECTION_STRING")

func init() {
	providers.Add("mssql", &MSSQL{})
}

// MSSQL is a migration provider for SQL Server.
type MSSQL struct{}

// GetAppliedMigrations queries the migration history table for all applied migrations.
func (*MSSQL) GetAppliedMigrations(ctx context.Context) ([]*migrations.Migration, error) {
	db, err := sql.Open("sqlserver", connectionString)
	if err != nil {
		return nil, err
	}

	err = ensureHistoryTable(ctx, db)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("SELECT [Id], [Name], [DateApplied] FROM [%s];", historyTableName)
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var appliedMigrations []*migrations.Migration

	for rows.Next() {
		var m migrations.Migration

		err := rows.Scan(
			&m.ID,
			&m.Name,
			&m.DateApplied,
		)
		if err != nil {
			return nil, err
		}

		appliedMigrations = append(appliedMigrations, &m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return appliedMigrations, nil
}

func ensureHistoryTable(ctx context.Context, db *sql.DB) error {
	query := fmt.Sprintf(
		`IF NOT EXISTS (SELECT [name] FROM sys.tables WHERE [name] = '%s')
		BEGIN
			CREATE TABLE [%s] (
				[Id] INT NOT NULL IDENTITY(1,1) PRIMARY KEY,
				[Name] VARCHAR(255) NOT NULL,
				[DateApplied] DATETIME NOT NULL
			);
		END`,
		historyTableName,
		historyTableName,
	)

	_, err := db.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// Apply applies the migration, m, to the database, as well as
// adding a record to the migration history table.
func (*MSSQL) Apply(ctx context.Context, m *migrations.Migration) error {
	db, err := sql.Open("sqlserver", connectionString)
	if err != nil {
		return err
	}

	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadUncommitted})
	if err != nil {
		return err
	}

	up, err := m.Up(ctx)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, up)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("INSERT INTO [%s] ([Name],[DateApplied]) VALUES (@name, GETUTCDATE());", historyTableName)
	_, err = tx.ExecContext(ctx, query, sql.Named("name", m.Name))
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

// Rollback rolls back the migration, m, then removed the
// record from the migration history table.
func (*MSSQL) Rollback(ctx context.Context, m *migrations.Migration) error {
	db, err := sql.Open("sqlserver", connectionString)
	if err != nil {
		return err
	}

	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadUncommitted})
	if err != nil {
		return err
	}

	down, err := m.Down(ctx)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, down)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("DELETE FROM [%s] WHERE [Name] = @name;", historyTableName)
	_, err = tx.ExecContext(ctx, query,
		sql.Named("name", m.Name))
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}
