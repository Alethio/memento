package commands

import (
	"database/sql"

	_ "github.com/Alethio/memento/migrations"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Manually run the database migrations",
	PreRun: func(cmd *cobra.Command, args []string) {
		bindViperToDBFlags(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		buildDBConnectionString()

		db, err := sql.Open("postgres", viper.GetString("db.connection-string"))
		if err != nil {
			log.Fatal(err)
		}

		err = goose.Up(db, "/tmp")
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	addDBFlags(migrateCmd)
}
