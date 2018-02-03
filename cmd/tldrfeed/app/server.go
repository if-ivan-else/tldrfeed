package app

import (
	"github.com/if-ivan-else/tldrfeed/internal/service"
	"github.com/spf13/cobra"
)

var port int
var iJSON bool

var config service.Config

func init() {
	serverCmd.PersistentFlags().IntVarP(&config.Port, "port", "p", 8080, "Port to bind to")
	serverCmd.PersistentFlags().BoolVarP(&config.IndentJSON, "indent-json", "i", false, "Indent JSON nicely in rendered API responses")

	RootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run tldrfeed server",
	Run:   runServer,
}

func runServer(cmd *cobra.Command, args []string) {
	// Start the tldr service
	service.Run(config)
}
