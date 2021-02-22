package migrations

import (
	"context"
)

// Provider is used to interface with a database system, such as SQL Server,
// to apply, rollback and get applied migrations.
type Provider interface {
	GetAppliedMigrations(ctx context.Context) ([]*Migration, error)
	Apply(ctx context.Context, name, content string) error
	Rollback(ctx context.Context, name, content string) error
}
