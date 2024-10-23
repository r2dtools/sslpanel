package migration

import (
	"backend/config"

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

			return m.Up()
		},
	}

	return &migrateUpCmd
}
