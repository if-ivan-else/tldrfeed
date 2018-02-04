package app

import (
	"log"

	"github.com/if-ivan-else/tldrfeed/internal/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var config service.Config

func init() {
	serverCmd.PersistentFlags().IntVarP(&config.Port, "port", "p", 8080, "Port to bind to")
	serverCmd.PersistentFlags().BoolVarP(&config.IndentJSON, "indent-json", "i", false, "Indent JSON nicely in rendered API responses")
	serverCmd.PersistentFlags().StringVarP(&config.DB, "db", "d", "0.0.0.0:27017/db", "DB connection URL")
	if err := viper.BindPFlag("db", serverCmd.PersistentFlags().Lookup("db")); err != nil {
		log.Fatal(err)
	}

	if err := viper.BindEnv("db", "DB_URL"); err != nil {
		log.Fatal(err)
	}

	RootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run tldrfeed server",
	Run:   runServer,
}

func runServer(cmd *cobra.Command, args []string) {
	// Start the tldr service
	config.DB = viper.GetString("db")
	service.Run(config)
}
