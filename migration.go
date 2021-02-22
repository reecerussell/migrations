package migrations

import (
	"time"
)

// Migration represents a migration.
type Migration struct {
	ID          int
	Name        string `yaml:"name"`
	DateApplied time.Time
	UpFile      string `yaml:"up"`
	DownFile    string `yaml:"down"`
}
