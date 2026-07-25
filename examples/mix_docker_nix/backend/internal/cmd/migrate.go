package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
)

func MigrateCmd(ctx context.Context) *cobra.Command {

	var filepath string
	var dburl string
	migrateCmd := &cobra.Command{
		Use:   "migrate up",
		Short: "Used for applying migration",
		RunE: func(cmd *cobra.Command, args []string) error {

			m, err := migrate.New(
				filepath,
				dburl,
			)
			if err != nil {
				return err
			}

			migrationErr := m.Up()
			sourceErr, databaseErr := m.Close()
			if errors.Is(migrationErr, migrate.ErrNoChange) {
				migrationErr = nil
			}
			if sourceErr != nil {
				sourceErr = fmt.Errorf("close migration source: %w", sourceErr)
			}
			if databaseErr != nil {
				databaseErr = fmt.Errorf("close migration database: %w", databaseErr)
			}
			return errors.Join(migrationErr, sourceErr, databaseErr)
		},
	}

	migrateCmd.PersistentFlags().StringVarP(&filepath, "filepath", "f", filepath, "Pass at start to give filepath of migration files.")
	migrateCmd.PersistentFlags().StringVarP(&dburl, "dburl", "u", dburl, "Pass at start to give dburl.")

	return migrateCmd
}
