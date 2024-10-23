package migration

import (
	"backend/config"

	"github.com/spf13/cobra"
)

var (
	steps uint
)

func GetMigrateDownCmd(config *config.Config) *cobra.Command {
	var migrateDownCmd = cobra.Command{
		Use:   "down",
		Short: "rollback all migrations from the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := getMigrationManager(config)

			if err != nil {
				return err
			}

			if steps != 0 {
				return m.Steps(-int(steps))
			}

			return m.Down()
		},
	}
	migrateDownCmd.PersistentFlags().UintVar(&steps, "steps", 0, "migrate down the specified number of steps")

	return &migrateDownCmd
}
