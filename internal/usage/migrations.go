package usage

import (
	"context"
	"fmt"
	"os"
)

func (s *Store) Migrate(ctx context.Context, migrationPath string) error {
	sql, err := os.ReadFile(migrationPath)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	_, err = s.db.Exec(ctx, string(sql))
	if err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	return nil
}
