package fixture

import (
	"backend/config"
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/spf13/cobra"
)

var name string
var supportedFixtureName map[string]bool = map[string]bool{
	"syntax": true,
}

func GetFixturesCmd(config *config.Config) *cobra.Command {
	var fixtureCmd = cobra.Command{
		Use:   "fixtures",
		Short: "Load fixtures to database",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := sql.Open(config.DbType, config.DbDsn)

			if err != nil {
				return err
			}

			directory := fmt.Sprintf("%s/fixtures/%s", config.BaseDir, config.Environment)
			var path func(*testfixtures.Loader) error

			if name == "" {
				path = testfixtures.Directory(directory)
			} else {
				if _, ok := supportedFixtureName[name]; !ok {
					return fmt.Errorf("fixture name '%s' is not supported", name)
				}

				path = testfixtures.Files(filepath.Join(directory, name+".yml"))
			}

			fixtures, err := testfixtures.New(
				testfixtures.Database(db),
				testfixtures.Dialect(config.DbType),
				path,
				testfixtures.DangerousSkipTestDatabaseCheck(),
			)

			if err != nil {
				return err
			}

			return fixtures.Load()
		},
	}

	fixtureCmd.PersistentFlags().StringVar(&name, "name", "", "fixture name")

	return &fixtureCmd
}
