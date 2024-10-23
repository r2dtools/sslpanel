package cli

import (
	"backend/config"
	"backend/internal/app/panel/adapters/cli/fixture"
	"backend/internal/app/panel/adapters/cli/migration"
	userStorage "backend/internal/app/panel/user/storage"
	"backend/internal/modules"
	"backend/internal/pkg/db"
	"backend/internal/pkg/logger"

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

	log, err := logger.NewLogger(config)

	if err != nil {
		return nil, err
	}

	database, err := db.GetDB(config)

	if err != nil {
		return nil, err
	}

	appUserStorage := userStorage.NewUserSqlStorage(database)

	migrationsCmd := migration.GetMigrationsCmd()

	rootCmd.AddCommand(migrationsCmd)
	migrationsCmd.AddCommand(migration.GetMigrateUpCmd(config))
	migrationsCmd.AddCommand(migration.GetMigrateDownCmd(config))
	migrationsCmd.AddCommand(migration.GetMigtateDropCmd(config))
	migrationsCmd.AddCommand(migration.GetMigrateForceCmd(config))

	rootCmd.AddCommand(fixture.GetFixturesCmd(config))

	modules.InitModulesCli(rootCmd, database, appUserStorage, log)

	return &App{
		cli: rootCmd,
	}, err
}
