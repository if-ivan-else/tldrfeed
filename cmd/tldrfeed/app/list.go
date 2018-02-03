package app

import (
	"os"

	"github.com/spf13/cobra"
)

func init() {

	listCmd.AddCommand(listUsersCmd)
	listCmd.AddCommand(listFeedsCmd)
	listCmd.AddCommand(listArticlesCmd)
	RootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List users, feeds or articles in tldrfeed",
	Run:   runList,
}

var listUsersCmd = &cobra.Command{
	Use:   "users",
	Short: "List users",
	Run:   runListUsers,
}

var listFeedsCmd = &cobra.Command{
	Use:   "feeds",
	Short: "List feeds",
	Run:   runListFeeds,
}

var listArticlesCmd = &cobra.Command{
	Use:   "articles",
	Short: "List articles",
	Run:   runListArticles,
}

func runList(cmd *cobra.Command, args []string) {
	cmd.Help()
	os.Exit(0)
}

func runListUsers(cmd *cobra.Command, args []string) {
	// TODO: Implement
}

func runListFeeds(cmd *cobra.Command, args []string) {
	// TODO: Implement
}

func runListArticles(cmd *cobra.Command, args []string) {
	// TODO: Implement
}
