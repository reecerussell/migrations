package mssql_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/reecerussell/migrations"
	"github.com/reecerussell/migrations/providers/mssql"
)

var testConnectionString = os.Getenv("MSSQL_CONNECTION_STRING")

func execute(db *sql.DB, queryf string, inlineArgs ...interface{}) {
	query := fmt.Sprintf(queryf, inlineArgs...)
	_, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
}

func TestGetAppliedMigrations_HavingOneAppliedMigration_ReturnsMigrationSuccessfully(t *testing.T) {
	db, err := sql.Open("sqlserver", testConnectionString)
	if err != nil {
		panic(err)
	}

	execute(db, `CREATE TABLE [__MigrationHistory] (
			[Id] INT NOT NULL IDENTITY(1,1) PRIMARY KEY,
			[Name] VARCHAR(255) NOT NULL,
			[DateApplied] DATETIME NOT NULL
		);`)

	execute(db, `INSERT INTO [__MigrationHistory] ([Name],[DateApplied]) 
		VALUES ('Test', GETUTCDATE())`)

	t.Cleanup(func() {
		execute(db, "DELETE FROM [__MigrationHistory]; DROP TABLE [__MigrationHistory];")
	})

	p := &mssql.MSSQL{ConnectionString: testConnectionString}
	appliedMigrations, err := p.GetAppliedMigrations(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, 1, len(appliedMigrations))
	assert.Equal(t, "Test", appliedMigrations[0].Name)
}

func TestGetAppliedMigrations_GivenInvalidConnectionString_ReturnsError(t *testing.T) {
	p := &mssql.MSSQL{}
	appliedMigrations, err := p.GetAppliedMigrations(context.TODO())
	assert.Nil(t, appliedMigrations)
	assert.NotNil(t, err)
}

func TestGetAppliedMigrations_WithInvalidHistoryTableStructure_ReturnsError(t *testing.T) {
	db, err := sql.Open("sqlserver", testConnectionString)
	if err != nil {
		panic(err)
	}

	// create migration history table to simulate a table with the same name,
	// but with a different table structure, i.e. without the "Name" column.
	execute(db, `CREATE TABLE [__MigrationHistory] (
			[Id] INT NOT NULL IDENTITY(1,1) PRIMARY KEY,
			[DateApplied] DATETIME NOT NULL
		);`)

	t.Cleanup(func() {
		execute(db, "DROP TABLE [__MigrationHistory];")
	})

	p := &mssql.MSSQL{ConnectionString: testConnectionString}
	appliedMigrations, err := p.GetAppliedMigrations(context.TODO())
	assert.Nil(t, appliedMigrations)
	assert.NotNil(t, err)
}

func TestGetAppliedMigrations_WithInvalidHistoryTableColumnTypes_ReturnsError(t *testing.T) {
	db, err := sql.Open("sqlserver", testConnectionString)
	if err != nil {
		panic(err)
	}

	// create migration history table to simulate a table with the same name,
	// but with a different column types, for example, "Id" being a string
	// instead of an int.
	execute(db, `CREATE TABLE [__MigrationHistory] (
			[Id] VARCHAR(10) NOT NULL PRIMARY KEY,
			[Name] VARCHAR(255) NOT NULL,
			[DateApplied] DATETIME NOT NULL
		);`)

	execute(db, `INSERT INTO [__MigrationHistory] ([Id],[Name],[DateApplied]) 
		VALUES ('001-1', 'Test', GETUTCDATE())`)

	t.Cleanup(func() {
		execute(db, "DELETE FROM [__MigrationHistory];")
		execute(db, "DROP TABLE [__MigrationHistory];")
	})

	p := &mssql.MSSQL{ConnectionString: testConnectionString}
	appliedMigrations, err := p.GetAppliedMigrations(context.TODO())
	assert.Nil(t, appliedMigrations)
	assert.NotNil(t, err)
}

func TestApply(t *testing.T) {
	db, err := sql.Open("sqlserver", testConnectionString)
	if err != nil {
		panic(err)
	}

	execute(db, `CREATE TABLE [__MigrationHistory] (
			[Id] INT NOT NULL IDENTITY(1,1) PRIMARY KEY,
			[Name] VARCHAR(255) NOT NULL,
			[DateApplied] DATETIME NOT NULL
		);`)

	file, err := os.Create("TestApply")
	if err != nil {
		panic(err)
	}

	file.Write([]byte(`CREATE TABLE [TestApply] (
		[Name] VARCHAR(255) NOT NULL
	)`))
	file.Close()

	t.Cleanup(func() {
		os.Remove("TestApply")

		execute(db, "DROP TABLE [TestApply];")
		execute(db, "DELETE FROM [__MigrationHistory];")
		execute(db, "DROP TABLE [__MigrationHistory];")
	})

	m := &migrations.Migration{
		Name:   "CreateTable",
		UpFile: "TestApply",
	}
	ctx := migrations.NewContext(context.TODO(), ".")

	p := &mssql.MSSQL{ConnectionString: testConnectionString}
	err = p.Apply(ctx, m)

	t.Run("Returns No Error", func(t *testing.T) {
		assert.NoError(t, err)
	})

	t.Run("Table Is Created", func(t *testing.T) {
		row := db.QueryRow("SELECT [name] FROM sys.tables WHERE [name] = 'TestApply'")
		var name string
		err = row.Scan(&name)

		assert.NoError(t, err)
		assert.Equal(t, "TestApply", name)
	})

	t.Run("Migration History Record Is Inserted", func(t *testing.T) {
		row := db.QueryRow("SELECT [Name] FROM [__MigrationHistory] WHERE [Name] = 'CreateTable'")
		var name string
		err = row.Scan(&name)

		assert.NoError(t, err)
		assert.Equal(t, "CreateTable", name)
	})
}

func TestApply_GivenMigrationWithInvalidSQL_ReturnsError(t *testing.T) {
	file, err := os.Create("TestApply_GivenMigrationWithInvalidSQL_ReturnsError")
	if err != nil {
		panic(err)
	}

	file.Write([]byte(`CREATE TABLE [TestApply] (
		[Name] VARCHAR(255) NO`)) // incomplete/invalid sql
	file.Close()

	t.Cleanup(func() {
		os.Remove("TestApply_GivenMigrationWithInvalidSQL_ReturnsError")
	})

	m := &migrations.Migration{
		Name:   "CreateTable",
		UpFile: "TestApply_GivenMigrationWithInvalidSQL_ReturnsError",
	}
	ctx := migrations.NewContext(context.TODO(), ".")

	p := &mssql.MSSQL{ConnectionString: testConnectionString}
	err = p.Apply(ctx, m)
	assert.NotNil(t, err)
}

func TestApply_GivenInvalidConnectionString_ReturnsError(t *testing.T) {
	p := &mssql.MSSQL{}
	err := p.Apply(context.TODO(), nil)
	assert.NotNil(t, err)
}

func TestApply_GivenMigrationWithMissingFile_ReturnError(t *testing.T) {
	m := &migrations.Migration{
		Name:   "CreateTable",
		UpFile: "TestApply_GivenMigrationWithMissingFile_ReturnError",
	}
	ctx := migrations.NewContext(context.TODO(), ".")

	p := &mssql.MSSQL{ConnectionString: testConnectionString}
	err := p.Apply(ctx, m)
	assert.True(t, os.IsNotExist(err))
}

func TestApply_WithInvalidHistoryTableStructure_ReturnError(t *testing.T) {
	db, err := sql.Open("sqlserver", testConnectionString)
	if err != nil {
		panic(err)
	}

	// Similarly to the other tests, this simulates a table already existing
	// in the database with the same name, with a different structure.
	// This misses the "Name" column, which will fail to insertion.
	execute(db, `CREATE TABLE [__MigrationHistory] (
			[Id] INT NOT NULL IDENTITY(1,1) PRIMARY KEY,
			[DateApplied] DATETIME NOT NULL
		);`)

	file, err := os.Create("TestApply")
	if err != nil {
		panic(err)
	}

	file.Write([]byte(`CREATE TABLE [TestApply] (
			[Name] VARCHAR(255) NOT NULL
		)`))
	file.Close()

	t.Cleanup(func() {
		os.Remove("TestApply")

		execute(db, "DROP TABLE [__MigrationHistory];")
	})

	m := &migrations.Migration{
		Name:   "CreateTable",
		UpFile: "TestApply",
	}
	ctx := migrations.NewContext(context.TODO(), ".")

	p := &mssql.MSSQL{ConnectionString: testConnectionString}
	err = p.Apply(ctx, m)
	assert.NotNil(t, err)
}

func TestRollback_GivenAppliedMigration_RollsBackSuucessful(t *testing.T) {
	db, err := sql.Open("sqlserver", testConnectionString)
	if err != nil {
		panic(err)
	}

	execute(db, `CREATE TABLE [__MigrationHistory] (
			[Id] INT NOT NULL IDENTITY(1,1) PRIMARY KEY,
			[Name] VARCHAR(255) NOT NULL,
			[DateApplied] DATETIME NOT NULL
		);`)

	execute(db, "INSERT INTO [__MigrationHistory] VALUES ('CreateTable', GETUTCDATE())")

	execute(db, `CREATE TABLE [TestRollback] (
			[Name] VARCHAR(255) NOT NULL
		)`)

	file, err := os.Create("TestRollback")
	if err != nil {
		panic(err)
	}

	file.Write([]byte(`DROP TABLE [TestRollback]`))
	file.Close()

	t.Cleanup(func() {
		os.Remove("TestRollback")

		execute(db, "DELETE FROM [__MigrationHistory];")
		execute(db, "DROP TABLE [__MigrationHistory];")
	})

	m := &migrations.Migration{
		Name:     "CreateTable",
		DownFile: "TestRollback",
	}
	ctx := migrations.NewContext(context.TODO(), ".")

	p := &mssql.MSSQL{ConnectionString: testConnectionString}
	err = p.Rollback(ctx, m)

	t.Run("Returns No Error", func(t *testing.T) {
		assert.NoError(t, err)
	})

	t.Run("Table Is Dropped", func(t *testing.T) {
		row := db.QueryRow("SELECT [name] FROM sys.tables WHERE [name] = 'TestRollback'")
		var name string
		err = row.Scan(&name)

		assert.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("Migration History Record Is Inserted", func(t *testing.T) {
		row := db.QueryRow("SELECT [Name] FROM [__MigrationHistory] WHERE [Name] = 'CreateTable'")
		var name string
		err = row.Scan(&name)

		assert.Equal(t, sql.ErrNoRows, err)
	})
}

func TestRollback_GivenInvalidConnectionString_ReturnsError(t *testing.T) {
	p := &mssql.MSSQL{}
	err := p.Rollback(context.TODO(), nil)
	assert.NotNil(t, err)
}

func TestRollback_GivenMigrationWithInvalidSQL_ReturnsError(t *testing.T) {
	file, err := os.Create("TestRollback")
	if err != nil {
		panic(err)
	}

	file.Write([]byte(`DROP TABLE [TestRollback`)) // incomplete/invalid sql
	file.Close()

	t.Cleanup(func() {
		os.Remove("TestRollback")
	})

	m := &migrations.Migration{
		Name:     "CreateTable",
		DownFile: "TestRollback",
	}
	ctx := migrations.NewContext(context.TODO(), ".")

	p := &mssql.MSSQL{ConnectionString: testConnectionString}
	err = p.Rollback(ctx, m)
	assert.NotNil(t, err)
}

func TestRollback_GivenMigrationWithMissingFile_ReturnError(t *testing.T) {
	m := &migrations.Migration{
		Name:     "CreateTable",
		DownFile: "TestRollback",
	}
	ctx := migrations.NewContext(context.TODO(), ".")

	p := &mssql.MSSQL{ConnectionString: testConnectionString}
	err := p.Rollback(ctx, m)
	assert.True(t, os.IsNotExist(err))
}

func TestRollback_WithInvalidHistoryTableStructure_ReturnError(t *testing.T) {
	db, err := sql.Open("sqlserver", testConnectionString)
	if err != nil {
		panic(err)
	}

	// Similarly to the other tests, this simulates a table already existing
	// in the database with the same name, with a different structure.
	// This misses the "Name" column, which will fail to delete.
	execute(db, `CREATE TABLE [__MigrationHistory] (
			[Id] INT NOT NULL IDENTITY(1,1) PRIMARY KEY,
			[DateApplied] DATETIME NOT NULL
		);`)

	execute(db, `CREATE TABLE [TestRollback] (
			[Name] VARCHAR(255) NOT NULL
		)`)

	file, err := os.Create("TestRollback")
	if err != nil {
		panic(err)
	}

	file.Write([]byte(`DROP TABLE [TestRollback];`))
	file.Close()

	t.Cleanup(func() {
		os.Remove("TestRollback")

		execute(db, "DROP TABLE [TestRollback];")
		execute(db, "DROP TABLE [__MigrationHistory];")
	})

	m := &migrations.Migration{
		Name:     "CreateTable",
		DownFile: "TestRollback",
	}
	ctx := migrations.NewContext(context.TODO(), ".")

	p := &mssql.MSSQL{ConnectionString: testConnectionString}
	err = p.Rollback(ctx, m)
	assert.NotNil(t, err)
}
