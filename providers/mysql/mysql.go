package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/reecerussell/migrations"
	"github.com/reecerussell/migrations/providers"

	// MySQL driver
	_ "github.com/go-sql-driver/mysql"
)

const defaultHistoryTableName = "__migration_history"

func init() {
	providers.Add("mysql", New)
}

// MySQL is a migration provider for MySQL.
type MySQL struct {
	ConnectionString string
	HistoryTableName string
	PrintStatements  bool
}

// New returns a new instance of MySQL. Implementing providers.ConstructorFunc,
// New takes a migrations.ConfigMap, which is used to populate HistoryTableName.
func New(conf migrations.ConfigMap) migrations.Provider {
	historyTableName := defaultHistoryTableName
	if v, _ := conf.String("historyTableName"); v != "" {
		historyTableName = v
	}
	printStatements := false
	if v, _ := conf.String("printStatements"); v == "true" {
		printStatements = true
	}
	return &MySQL{
		ConnectionString: os.Getenv("CONNECTION_STRING"),
		HistoryTableName: historyTableName,
		PrintStatements:  printStatements,
	}
}

// GetAppliedMigrations queries the migration history table for all applied migrations.
func (p *MySQL) GetAppliedMigrations(ctx context.Context) ([]*migrations.Migration, error) {
	db, err := p.openConn(ctx)
	if err != nil {
		return nil, err
	}
	p.ensureHistoryTable(ctx, db)
	query := fmt.Sprintf("SELECT `id`, `name`, `date_applied` FROM `%s`;", p.HistoryTableName)
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
func (p *MySQL) ensureHistoryTable(ctx context.Context, db *sql.DB) {
	query := fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS `%s` ("+
			"`id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY,"+
			"`name` VARCHAR(255) NOT NULL,"+
			"`date_applied` DATETIME NOT NULL"+
			");",
		p.HistoryTableName,
	)

	// This should never return an error, as it's given a valid
	// *sql.DB instance, with an open connection. Plus the query
	// is valid.
	db.ExecContext(ctx, query)
}

// Apply applies the migration, m, to the database, as well as
// adding a record to the migration history table.
func (p *MySQL) Apply(ctx context.Context, name, content string) error {
	db, err := p.openConn(ctx)
	if err != nil {
		return err
	}
	tx, _ := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadUncommitted})
	statements := strings.Split(content, ";")
	for _, statement := range statements {
		if strings.TrimSpace(statement) == "" {
			continue
		}
		if p.PrintStatements {
			fmt.Printf("Executing the following statement:\n%s\n", statement)
		}
		_, err = tx.ExecContext(ctx, statement)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	query := fmt.Sprintf("INSERT INTO `%s` (`name`,`date_applied`) VALUES (?, UTC_TIMESTAMP());", p.HistoryTableName)
	_, err = tx.ExecContext(ctx, query, name)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// Rollback rolls back the migration, m, then removed the
// record from the migration history table.
func (p *MySQL) Rollback(ctx context.Context, name, content string) error {
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
	query := fmt.Sprintf("DELETE FROM `%s` WHERE `name` = ?;", p.HistoryTableName)
	_, err = tx.ExecContext(ctx, query, name)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (p *MySQL) openConn(ctx context.Context) (*sql.DB, error) {
	db, _ := sql.Open("mysql", p.ConnectionString)
	err := db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
