package cli

import (
	"backend/config"
	"backend/internal/app/panel/adapters/cli/fixture"
	"backend/internal/app/panel/adapters/cli/migration"

	"github.com/spf13/cobra"
)

type App struct {
	cli *cobra.Command
}

func (app *App) Run() error {
	return app.cli.Execute()
}

func GetApp() (*App, error) {
	rootCmd := &cobra.Command{
		Use:   "r2cli",
		Short: "R2 Control Panel CLI",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}

	config, err := config.GetConfig()

	if err != nil {
		return nil, err
	}

	migrationsCmd := migration.GetMigrationsCmd()

	rootCmd.AddCommand(migrationsCmd)
	migrationsCmd.AddCommand(migration.GetMigrateUpCmd(config))
	migrationsCmd.AddCommand(migration.GetMigrateDownCmd(config))
	migrationsCmd.AddCommand(migration.GetMigtateDropCmd(config))
	migrationsCmd.AddCommand(migration.GetMigrateForceCmd(config))

	rootCmd.AddCommand(fixture.GetFixturesCmd(config))

	return &App{
		cli: rootCmd,
	}, err
}
