package migrations

import (
	"context"
	"sync"
)

// In-memory store of registered migration providers.
var (
	mu   sync.Mutex
	prvs []Provider
)

// Provider is used to interface with a database system, such as SQL Server,
// to apply, rollback and get applied migrations.
type Provider interface {
	GetAppliedMigrations(ctx context.Context) ([]*Migration, error)
	Apply(ctx context.Context, m *Migration) error
	Rollback(ctx context.Context, m *Migration) error
}
