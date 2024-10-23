package migration

import (
	"backend/config"
	"errors"

	"github.com/spf13/cobra"
)

var version uint

func GetMigrateForceCmd(config *config.Config) *cobra.Command {
	var migrateForceCmd = cobra.Command{
		Use:   "force",
		Short: "set migration version",
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := getMigrationManager(config)

			if err != nil {
				return err
			}

			if version == 0 {
				return errors.New("invalid version value")
			}

			return m.Force(int(version))
		},
	}

	migrateForceCmd.PersistentFlags().UintVar(&version, "version", 0, "migration version")

	return &migrateForceCmd
}
