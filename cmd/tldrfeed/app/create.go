package app

import (
	"os"

	"github.com/spf13/cobra"
)

func init() {

	createCmd.AddCommand(createUserCmd)
	createCmd.AddCommand(createFeedCmd)
	createCmd.AddCommand(createArticleCmd)
	RootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create users, feeds or articles in tldrfeed",
	Run:   runCreate,
}

var createUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Create new user",
	Run:   runCreateUser,
}

var createFeedCmd = &cobra.Command{
	Use:   "feed",
	Short: "Create new feed",
	Run:   runCreateFeed,
}

var createArticleCmd = &cobra.Command{
	Use:   "article",
	Short: "Create new article",
	Run:   runCreateArticle,
}

func runCreate(cmd *cobra.Command, args []string) {
	cmd.Help()
	os.Exit(0)
}

func runCreateUser(cmd *cobra.Command, args []string) {

}

func runCreateFeed(cmd *cobra.Command, args []string) {

}

func runCreateArticle(cmd *cobra.Command, args []string) {

}
