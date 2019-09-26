package commands

import (
	"time"

	_ "git.aleth.io/alethio/memento/migrations"
	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var queueCmd = &cobra.Command{
	Use:   "queue",
	Short: "Manually add a block to the todo queue",
	PreRun: func(cmd *cobra.Command, args []string) {
		bindViperToRedisFlags(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		r := redis.NewClient(&redis.Options{
			Addr:        viper.GetString("redis.server"),
			Password:    viper.GetString("REDIS_PASSWORD"),
			DB:          0,
			ReadTimeout: time.Second * 1,
		})

		err := r.Ping().Err()
		if err != nil {
			log.Fatal(err)
			return
		}

		list := viper.GetString("redis.list")

		block := viper.GetInt64("block")
		if block > 0 {
			err := addTodo(r, list, block)
			if err != nil {
				log.Fatal(err)
			}
			return
		}

		from := viper.GetInt64("from")
		to := viper.GetInt64("to")
		if from > 0 && to > 0 {
			for i := from; i <= to; i++ {
				err := addTodo(r, list, i)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	},
}

func addTodo(r *redis.Client, list string, number int64) error {
	log.WithField("block", number).Info("adding block to todo")
	return r.ZAdd(list, redis.Z{
		Score:  float64(number),
		Member: number,
	}).Err()
}

func init() {
	addRedisFlags(queueCmd)

	queueCmd.Flags().Int64("block", -1, "Add a single block in the todo queue")
	viper.BindPFlag("block", queueCmd.Flag("block"))

	queueCmd.Flags().Int64("from", -1, "Add a series of blocks into the todo queue starting from the provided number (only use in combination with --to)")
	viper.BindPFlag("from", queueCmd.Flag("from"))

	queueCmd.Flags().Int64("to", -1, "Add a series of blocks into the todo queue ending with the provided number, inclusive (only use in combination with --from)")
	viper.BindPFlag("to", queueCmd.Flag("to"))
}
