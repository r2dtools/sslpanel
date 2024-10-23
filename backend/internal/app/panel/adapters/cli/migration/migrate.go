package migration

import (
	"backend/config"
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
)

func GetMigrationsCmd() *cobra.Command {
	var migrationsCmd = cobra.Command{
		Use:   "migrations",
		Short: "Manage database migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}

	return &migrationsCmd
}

func getMigrationManager(config *config.Config) (*migrate.Migrate, error) {
	db, err := sql.Open(config.DbType, config.DbDsn)

	if err != nil {
		return nil, err
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{})

	if err != nil {
		return nil, err
	}

	migrate, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s/migrations", config.BaseDir),
		config.DbType,
		driver,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create migration: %v", err)
	}

	return migrate, nil
}
