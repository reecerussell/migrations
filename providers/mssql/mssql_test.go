package mssql

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func execute(db *sql.DB, queryf string, inlineArgs ...interface{}) {
	query := fmt.Sprintf(queryf, inlineArgs...)
	_, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
}

func TestGetAppliedMigrations_HavingOneAppliedMigration_ReturnsMigrationSuccessfully(t *testing.T) {
	db, err := sql.Open("sqlserver", connectionString)
	if err != nil {
		panic(err)
	}

	execute(db, `IF NOT EXISTS (SELECT [name] FROM sys.tables WHERE [name] = '__MigrationHistory')
		BEGIN
			CREATE TABLE [__MigrationHistory] (
				[Id] INT NOT NULL IDENTITY(1,1) PRIMARY KEY,
				[Name] VARCHAR(255) NOT NULL,
				[DateApplied] DATETIME NOT NULL
			);
		END`)

	execute(db, `INSERT INTO [__MigrationHistory] ([Name],[DateApplied]) 
		VALUES ('Test', GETUTCDATE())`)

	t.Cleanup(func() {
		execute(db, "DELETE FROM [__MigrationHistory]; DROP TABLE [__MigrationHistory];")
	})

	p := &MSSQL{}
	appliedMigrations, err := p.GetAppliedMigrations(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, 1, len(appliedMigrations))
	assert.Equal(t, "Test", appliedMigrations[0].Name)
}
