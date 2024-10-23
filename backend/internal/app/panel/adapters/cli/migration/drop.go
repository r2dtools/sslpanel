package migration

import (
	"backend/config"

	"github.com/spf13/cobra"
)

func GetMigtateDropCmd(config *config.Config) *cobra.Command {
	var migrateDropCmd = cobra.Command{
		Use:   "drop",
		Short: "drop all data from the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := getMigrationManager(config)

			if err != nil {
				return err
			}

			return m.Drop()
		},
	}

	return &migrateDropCmd
}
