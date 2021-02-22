package migrations

import (
	"context"
	"fmt"
)

// Rollback rolls back all applied migrations, up to the target (if any), using the given provider, p.
func Rollback(ctx context.Context, cm []*Migration, p Provider, fr FileReader, targetName string) error {
	am, err := p.GetAppliedMigrations(ctx)
	if err != nil {
		return err
	}

	for i := len(cm) - 1; i >= 0; i-- {
		m := cm[i]

		fmt.Printf("Rolling back %s...\t", m.Name)

		if !isApplied(am, m.Name) {
			fmt.Printf("skipping.\n")
			continue
		}

		content, err := fr.Read(m.DownFile)
		if err != nil {
			fmt.Printf("\nFailed to read migration file: %s.\n", m.DownFile)
			return err
		}

		err = p.Rollback(ctx, m.Name, content)
		if err != nil {
			fmt.Printf("\nFailed to rollback migration %s.\n", m.Name)

			return err
		}

		fmt.Printf("done.\n")

		if targetName != "" && targetName == m.Name {
			return nil
		}
	}

	return nil
}
