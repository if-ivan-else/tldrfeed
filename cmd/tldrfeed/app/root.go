package app

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

// RootCmd is the root command that is executed when sonobuoy is run without
// any subcommands.
var RootCmd = &cobra.Command{
	Use:   "tldrfeed",
	Short: "Easy to use JSON news feeds for everyone.",
	Long:  "tldrfeed is a simple JSON news feed subscription service",
	Run:   rootCmd,
}

func rootCmd(cmd *cobra.Command, args []string) {
	// Do nothing when not given a subcommand
	err := cmd.Help()
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
