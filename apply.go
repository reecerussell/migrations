package migrations

import (
	"context"
	"fmt"
)

// Apply applies all unapplied migrations, up to the target (if any), using the given provider, p.
func Apply(ctx context.Context, cm []*Migration, p Provider, targetName string) error {
	am, err := p.GetAppliedMigrations(ctx)
	if err != nil {
		return err
	}

	for _, m := range cm {
		fmt.Printf("Applying %s...\t", m.Name)

		if isApplied(am, m.Name) {
			fmt.Printf("skipping.\n")
			continue
		}

		err := p.Apply(ctx, m)
		if err != nil {
			fmt.Printf("\nFailed to apply migration %s.\n", m.Name)

			return err
		}

		fmt.Printf("done.\n")

		if targetName != "" && targetName == m.Name {
			return nil
		}
	}

	return nil
}

func isApplied(applied []*Migration, name string) bool {
	for _, m := range applied {
		if m.Name == name {
			return true
		}
	}

	return false
}
