package migrations

import (
	"context"
)

// Provider is used to interface with a database system, such as SQL Server,
// to apply, rollback and get applied migrations.
type Provider interface {
	// GetAppliedMigrations returns an array of all Migrations that
	// have already been applied. If an error is returned it would
	// be specific to the provider in use.
	GetAppliedMigrations(ctx context.Context) ([]*Migration, error)

	// Apply applies a migration with the given name, using the content
	// provided. A record should be held of the migration application.
	Apply(ctx context.Context, name, content string) error

	// Rollback reverts an already-applied migration, with the given
	// name and content. If successful, the record of the migration
	// should be removed, to prevent any issues with other migrations.
	Rollback(ctx context.Context, name, content string) error
}
