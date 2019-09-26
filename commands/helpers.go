package commands

import (
	"fmt"

	"github.com/gin-gonic/gin"
	formatter "github.com/kwix/logrus-module-formatter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func initLogging() {
	if verbose && logging == "" {
		logging = "*=debug"
	}

	if vverbose && logging == "" {
		logging = "*=trace"
	}

	if logging == "" {
		logging = "*=info"
	}

	gin.SetMode(gin.DebugMode)

	modules := formatter.NewModulesMap(logging)
	if level, exists := modules["gin"]; exists {
		if level < logrus.DebugLevel {
			gin.SetMode(gin.ReleaseMode)
		}
	} else {
		level := modules["*"]
		if level < logrus.DebugLevel {
			gin.SetMode(gin.ReleaseMode)
		}
	}

	f, err := formatter.New(modules)
	if err != nil {
		panic(err)
	}

	logrus.SetFormatter(f)

	log.Debug("Debug mode")
}

func addDBFlags(cmd *cobra.Command) {
	cmd.Flags().String("db.connection-string", "", "Postgres connection string.")
	cmd.Flags().String("db.host", "localhost", "Database host")
	cmd.Flags().String("db.port", "5432", "Database port")
	cmd.Flags().String("db.sslmode", "disable", "Database sslmode")
	cmd.Flags().String("db.dbname", "coriolis", "Database name")
	cmd.Flags().String("db.user", "", "Database user (also allowed via PG_USER env)")
}

func bindViperToDBFlags(cmd *cobra.Command) {
	viper.BindPFlag("db.connection-string", cmd.Flag("db.connection-string"))
	viper.BindPFlag("db.host", cmd.Flag("db.host"))
	viper.BindPFlag("db.port", cmd.Flag("db.port"))
	viper.BindPFlag("db.sslmode", cmd.Flag("db.sslmode"))
	viper.BindPFlag("db.dbname", cmd.Flag("db.dbname"))
	viper.BindPFlag("db.user", cmd.Flag("db.user"))
}

func addRedisFlags(cmd *cobra.Command) {
	cmd.Flags().String("redis.server", "localhost:6379", "Redis server URL")
	cmd.Flags().String("redis.list", "todo", "The name of the list to be used for task management")
}

func bindViperToRedisFlags(cmd *cobra.Command) {
	viper.BindPFlag("redis.server", cmd.Flag("redis.server"))
	viper.BindPFlag("redis.list", cmd.Flag("redis.list"))
}

func buildDBConnectionString() {
	if viper.GetString("db.connection-string") == "" {
		var user, pass string
		if !viper.IsSet("db.user") {
			user = viper.GetString("PG_USER")
		} else {
			user = viper.GetString("db.user")
		}

		pass = viper.GetString("PG_PASSWORD")

		p := fmt.Sprintf("host=%s port=%s sslmode=%s dbname=%s user=%s password=%s", viper.GetString("db.host"), viper.GetString("db.port"), viper.GetString("db.sslmode"), viper.GetString("db.dbname"), user, pass)
		viper.Set("db.connection-string", p)
	}
}
