package app

import (
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/if-ivan-else/tldrfeed/api"
	"github.com/spf13/cobra"
)

var userID string

func init() {
	listCmd.PersistentFlags().StringVarP(&url, "url", "u", "http://localhost:8080", "tldrfeed service URL")

	listCmd.AddCommand(listUsersCmd)
	listCmd.AddCommand(listFeedsCmd)

	listArticlesCmd.PersistentFlags().StringVarP(&feedID, "feed", "f", "", "Feed ID")
	listArticlesCmd.PersistentFlags().StringVar(&userID, "user", "", "User ID")
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
	c := api.NewClient(url)
	users, err := c.ListUsers()
	if err != nil {
		log.Fatalf("Failed to list Users: %s", err)
	}
	log.Print("Users:\n")
	for _, u := range users {
		spew.Printf("%+v\n", u)
	}
}

func runListFeeds(cmd *cobra.Command, args []string) {
	c := api.NewClient(url)
	feeds, err := c.ListFeeds()
	if err != nil {
		log.Fatalf("Failed to list Feeds: %s", err)
	}
	log.Print("Feeds:\n")
	for _, f := range feeds {
		spew.Printf("%+v\n", f)
	}
}

func runListArticles(cmd *cobra.Command, args []string) {
	c := api.NewClient(url)
	articles := []api.Article{}
	var err error
	if userID == "" {
		log.Printf("Articles in feed %s:", feedID)
		articles, err = c.ListArticles(feedID)
	} else {
		if feedID == "" {
			log.Printf("Articles for user %s in all Feeds", userID)
		} else {
			log.Printf("Articles for user %s in Feed %s):", userID, feedID)
		}
		articles, err = c.ListUserArticles(userID, feedID)
	}
	if err != nil {
		log.Fatalf("Failed to list Articles: %s", err)
	}
	for _, a := range articles {
		spew.Printf("%+v\n", a)
	}
}
