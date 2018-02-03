package app

import (
	"fmt"

	"github.com/if-ivan-else/tldrfeed/internal/buildinfo"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print tldrfeed version",
	Run:   runVersion,
}

func runVersion(cmd *cobra.Command, args []string) {
	fmt.Println(buildinfo.Version)
}
