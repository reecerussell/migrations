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

func init() {
	providers.Add("mssql", &MSSQL{
		ConnectionString: os.Getenv("CONNECTION_STRING"),
	})
}

// MSSQL is a migration provider for SQL Server.
type MSSQL struct {
	ConnectionString string
}

// GetAppliedMigrations queries the migration history table for all applied migrations.
func (p *MSSQL) GetAppliedMigrations(ctx context.Context) ([]*migrations.Migration, error) {
	db, err := p.openConn(ctx)
	if err != nil {
		return nil, err
	}

	ensureHistoryTable(ctx, db)

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

	return appliedMigrations, nil
}

// ensureHistoryTable ensures the table with the name historyTableName exists.
// Should be provided a valid instance of *sql.DB.
func ensureHistoryTable(ctx context.Context, db *sql.DB) {
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

	// This should never return an error, as it's given a valid
	// *sql.DB instance, with an open connection. Plus the query
	// is valid.
	db.ExecContext(ctx, query)
}

// Apply applies the migration, m, to the database, as well as
// adding a record to the migration history table.
func (p *MSSQL) Apply(ctx context.Context, name, content string) error {
	db, err := p.openConn(ctx)
	if err != nil {
		return err
	}

	tx, _ := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadUncommitted})
	_, err = tx.ExecContext(ctx, content)
	if err != nil {
		tx.Rollback()
		return err
	}

	query := fmt.Sprintf("INSERT INTO [%s] ([Name],[DateApplied]) VALUES (@name, GETUTCDATE());", historyTableName)
	_, err = tx.ExecContext(ctx, query, sql.Named("name", name))
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

// Rollback rolls back the migration, m, then removed the
// record from the migration history table.
func (p *MSSQL) Rollback(ctx context.Context, name, content string) error {
	var err error

	db, err := p.openConn(ctx)
	if err != nil {
		return err
	}

	tx, _ := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadUncommitted})
	_, err = tx.ExecContext(ctx, content)
	if err != nil {
		tx.Rollback()
		return err
	}

	query := fmt.Sprintf("DELETE FROM [%s] WHERE [Name] = @name;", historyTableName)
	_, err = tx.ExecContext(ctx, query, sql.Named("name", name))
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (p *MSSQL) openConn(ctx context.Context) (*sql.DB, error) {
	db, _ := sql.Open("sqlserver", p.ConnectionString)
	err := db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
