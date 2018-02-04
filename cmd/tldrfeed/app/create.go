package app

import (
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"

	"github.com/if-ivan-else/tldrfeed/api"
	"github.com/spf13/cobra"
)

var url string
var name string
var title string
var body string
var feedID string

func init() {

	createCmd.PersistentFlags().StringVarP(&url, "url", "u", "http://localhost:8080", "tldrfeed service URL")

	createUserCmd.PersistentFlags().StringVarP(&name, "name", "n", "", "User name")
	createCmd.AddCommand(createUserCmd)

	createFeedCmd.PersistentFlags().StringVarP(&feedID, "feed", "f", "", "Feed ID")
	createFeedCmd.PersistentFlags().StringVarP(&name, "name", "n", "", "Feed name")
	createCmd.AddCommand(createFeedCmd)

	createArticleCmd.PersistentFlags().StringVarP(&title, "title", "t", "", "Article title")
	createArticleCmd.PersistentFlags().StringVarP(&body, "body", "b", "", "Article body")
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
	_ = cmd.Help()
	os.Exit(0)
}

func runCreateUser(cmd *cobra.Command, args []string) {
	c := api.NewClient(url)
	u, err := c.CreateUser(name)
	if err != nil {
		log.Fatalf("Failed to create User: %s", err.Error())
	}
	spew.Printf("User created: %v", u)
}

func runCreateFeed(cmd *cobra.Command, args []string) {
	c := api.NewClient(url)
	f, err := c.CreateFeed(name)
	if err != nil {
		log.Fatalf("Failed to create Feed: %s", err.Error())
	}
	spew.Printf("Feed created: %v", f)
}

func runCreateArticle(cmd *cobra.Command, args []string) {
	c := api.NewClient(url)
	f, err := c.CreateArticle(feedID, title, body)
	if err != nil {
		log.Fatalf("Failed to create Feed: %s", err.Error())
	}
	spew.Printf("Feed created: %v", f)
}
