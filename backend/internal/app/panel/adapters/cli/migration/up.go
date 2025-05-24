package migration

import (
	"backend/config"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/spf13/cobra"
)

func GetMigrateUpCmd(config *config.Config) *cobra.Command {
	var migrateUpCmd = cobra.Command{
		Use:   "up",
		Short: "apply all migrations to the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := getMigrationManager(config)

			if err != nil {
				return err
			}

			err = m.Up()

			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("migrations: no changes")

				return nil
			}

			return err
		},
	}

	return &migrateUpCmd
}
