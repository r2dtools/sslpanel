package cmd

import (
	"backend/db"
	"backend/db/seeds"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	// initialize migration from a file
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate up|down|drop",
	Short: "migrates application database",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		instance := db.GetDB().DB()
		driver, err := postgres.WithInstance(instance, &postgres.Config{})

		if err != nil {
			return err
		}

		m, err := migrate.NewWithDatabaseInstance("file://db/migrations", "postgres", driver)

		if err != nil {
			return err
		}

		currentVersion, _, _ := m.Version()

		switch direction := args[0]; direction {
		case "up":
			err = m.Up()

			if err != nil {
				break
			}

			version, _, _ := m.Version()
			err = seeds.ApplyAll(instance, version)

			if err != nil {
				m.Migrate(currentVersion)
			}
		case "down":
			err = m.Down()
		case "drop":
			err = m.Drop()
		}

		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
