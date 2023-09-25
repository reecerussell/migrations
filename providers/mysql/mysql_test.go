package mysql_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/reecerussell/migrations"
	"github.com/reecerussell/migrations/providers/mysql"
)

var testConnectionString = os.Getenv("MYSQL_CONNECTION_STRING")

func execute(db *sql.DB, queryf string, inlineArgs ...interface{}) {
	query := fmt.Sprintf(queryf, inlineArgs...)
	_, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
}

func TestNew(t *testing.T) {
	cnf := migrations.ConfigMap{
		"historyTableName": "MyMigrationsTable",
	}
	p := mysql.New(cnf).(*mysql.MySQL)

	assert.Equal(t, "MyMigrationsTable", p.HistoryTableName)
}

func TestGetAppliedMigrations_HavingOneAppliedMigration_ReturnsMigrationSuccessfully(t *testing.T) {
	db, err := sql.Open("mysql", testConnectionString)
	if err != nil {
		panic(err)
	}

	execute(db, `CREATE TABLE __MigrationHistory (
			id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			date_applied DATETIME NOT NULL
		);`)

	execute(db, `INSERT INTO __MigrationHistory (name,date_applied) 
		VALUES ('Test', UTC_TIMESTAMP())`)

	t.Cleanup(func() {
		execute(db, "DROP TABLE __MigrationHistory;")
	})

	p := &mysql.MySQL{
		ConnectionString: testConnectionString,
		HistoryTableName: "__MigrationHistory",
	}
	appliedMigrations, err := p.GetAppliedMigrations(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, 1, len(appliedMigrations))
	assert.Equal(t, "Test", appliedMigrations[0].Name)
}

func TestGetAppliedMigrations_GivenInvalidConnectionString_ReturnsError(t *testing.T) {
	p := &mysql.MySQL{}
	appliedMigrations, err := p.GetAppliedMigrations(context.TODO())
	assert.Nil(t, appliedMigrations)
	assert.NotNil(t, err)
}

func TestGetAppliedMigrations_WithInvalidHistoryTableStructure_ReturnsError(t *testing.T) {
	db, err := sql.Open("mysql", testConnectionString)
	if err != nil {
		panic(err)
	}

	// create migration history table to simulate a table with the same name,
	// but with a different table structure, i.e. without the "Name" column.
	execute(db, `CREATE TABLE __MigrationHistory (
			id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
			date_applied DATETIME NOT NULL
		);`)

	t.Cleanup(func() {
		execute(db, "DROP TABLE __MigrationHistory;")
	})

	p := &mysql.MySQL{
		ConnectionString: testConnectionString,
		HistoryTableName: "__MigrationHistory",
	}
	appliedMigrations, err := p.GetAppliedMigrations(context.TODO())
	assert.Nil(t, appliedMigrations)
	assert.NotNil(t, err)
}

func TestGetAppliedMigrations_WithInvalidHistoryTableColumnTypes_ReturnsError(t *testing.T) {
	db, err := sql.Open("mysql", testConnectionString)
	if err != nil {
		panic(err)
	}

	// create migration history table to simulate a table with the same name,
	// but with a different column types, for example, "Id" being a string
	// instead of an int.
	execute(db, `CREATE TABLE __MigrationHistory (
			id VARCHAR(10) NOT NULL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			date_applied DATETIME NOT NULL
		);`)

	execute(db, `INSERT INTO __MigrationHistory (id, name, date_applied) 
		VALUES ('001-1', 'Test', UTC_TIMESTAMP())`)

	t.Cleanup(func() {
		execute(db, "DELETE FROM __MigrationHistory;")
		execute(db, "DROP TABLE __MigrationHistory;")
	})

	p := &mysql.MySQL{
		ConnectionString: testConnectionString,
		HistoryTableName: "__MigrationHistory",
	}
	appliedMigrations, err := p.GetAppliedMigrations(context.TODO())
	assert.Nil(t, appliedMigrations)
	assert.NotNil(t, err)
}

func TestApply(t *testing.T) {
	db, err := sql.Open("mysql", testConnectionString)
	if err != nil {
		panic(err)
	}

	execute(db, `CREATE TABLE __MigrationHistory (
			id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			date_applied DATETIME NOT NULL
		);`)

	t.Cleanup(func() {
		execute(db, "DROP TABLE TestApply;")
		execute(db, "DELETE FROM __MigrationHistory;")
		execute(db, "DROP TABLE __MigrationHistory;")
	})

	p := &mysql.MySQL{
		ConnectionString: testConnectionString,
		HistoryTableName: "__MigrationHistory",
	}
	err = p.Apply(context.TODO(), "CreateTable", `CREATE TABLE TestApply (
			name VARCHAR(255) NOT NULL
		)`)

	t.Run("Returns No Error", func(t *testing.T) {
		assert.NoError(t, err)
	})

	t.Run("Table Is Created", func(t *testing.T) {
		row := db.QueryRow("SELECT table_name FROM information_schema.tables WHERE table_name = 'TestApply'")
		var name string
		err = row.Scan(&name)

		assert.NoError(t, err)
		assert.Equal(t, "TestApply", name)
	})

	t.Run("Migration History Record Is Inserted", func(t *testing.T) {
		row := db.QueryRow("SELECT name FROM __MigrationHistory WHERE name = 'CreateTable'")
		var name string
		err = row.Scan(&name)

		assert.NoError(t, err)
		assert.Equal(t, "CreateTable", name)
	})
}

func TestApply_GivenMigrationWithInvalidSQL_ReturnsError(t *testing.T) {
	p := &mysql.MySQL{
		ConnectionString: testConnectionString,
		HistoryTableName: "__MigrationHistory",
	}
	err := p.Apply(context.TODO(), "TestApply", `CREATE TABLE TestApply (
		name VARCHAR(255) NO`)
	assert.NotNil(t, err)
}

func TestApply_GivenInvalidConnectionString_ReturnsError(t *testing.T) {
	p := &mysql.MySQL{}
	err := p.Apply(context.TODO(), "", "")
	assert.NotNil(t, err)
}

func TestApply_WithInvalidHistoryTableStructure_ReturnError(t *testing.T) {
	db, err := sql.Open("mysql", testConnectionString)
	if err != nil {
		panic(err)
	}

	// Similarly to the other tests, this simulates a table already existing
	// in the database with the same name, with a different structure.
	// This misses the "Name" column, which will fail to insertion.
	execute(db, `CREATE TABLE __MigrationHistory (
			id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
			date_applied DATETIME NOT NULL
		);`)

	t.Cleanup(func() {
		execute(db, "DROP TABLE __MigrationHistory;")
	})

	p := &mysql.MySQL{
		ConnectionString: testConnectionString,
		HistoryTableName: "__MigrationHistory",
	}
	err = p.Apply(context.TODO(), "TestApply", `CREATE TABLE TestApply (
			name VARCHAR(255) NOT NULL
		)`)
	assert.NotNil(t, err)
}

func TestRollback_GivenAppliedMigration_RollsBackSuucessful(t *testing.T) {
	db, err := sql.Open("mysql", testConnectionString)
	if err != nil {
		panic(err)
	}

	execute(db, `CREATE TABLE __MigrationHistory (
			id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			date_applied DATETIME NOT NULL
		);`)

	execute(db, "INSERT INTO __MigrationHistory VALUES (1, 'CreateTable', UTC_TIMESTAMP())")

	execute(db, `CREATE TABLE TestRollback (
			name VARCHAR(255) NOT NULL
		)`)

	t.Cleanup(func() {
		execute(db, "DELETE FROM __MigrationHistory;")
		execute(db, "DROP TABLE __MigrationHistory;")
	})

	p := &mysql.MySQL{
		ConnectionString: testConnectionString,
		HistoryTableName: "__MigrationHistory",
	}
	err = p.Rollback(context.TODO(), "CreateTable", `DROP TABLE TestRollback`)

	t.Run("Returns No Error", func(t *testing.T) {
		assert.NoError(t, err)
	})

	t.Run("Table Is Dropped", func(t *testing.T) {
		row := db.QueryRow("SELECT table_name FROM information_schema.tables WHERE table_name = 'TestRollback'")
		var name string
		err = row.Scan(&name)

		assert.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("Migration History Record Is Inserted", func(t *testing.T) {
		row := db.QueryRow("SELECT name FROM __MigrationHistory WHERE name = 'CreateTable'")
		var name string
		err = row.Scan(&name)

		assert.Equal(t, sql.ErrNoRows, err)
	})
}

func TestRollback_GivenInvalidConnectionString_ReturnsError(t *testing.T) {
	p := &mysql.MySQL{}
	err := p.Rollback(context.TODO(), "", "")
	assert.NotNil(t, err)
}

func TestRollback_GivenMigrationWithInvalidSQL_ReturnsError(t *testing.T) {
	p := &mysql.MySQL{
		ConnectionString: testConnectionString,
		HistoryTableName: "__MigrationHistory",
	}
	err := p.Rollback(context.TODO(), "CreateTable", `DROP TABLE TestRollback'`) // invalid sql
	assert.NotNil(t, err)
}
