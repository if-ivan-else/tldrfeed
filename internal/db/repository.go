package db

import (
	"github.com/if-ivan-else/tldrfeed/api"
)

// Repository defines an interface with persistence layer for Users, Feeds, and Articles entities
type Repository interface {
	CreateUser(name string) (*api.User, error)

	ListUsers() ([]api.User, error)

	GetUser(userID string) (*api.User, error)

	CreateFeed(name string) (*api.Feed, error)

	ListFeeds() ([]api.Feed, error)

	GetFeed(feedID string) (*api.Feed, error)

	ListFeedArticles(feedID string) ([]api.Article, error)

	CreateFeedArticle(feedID string, articleTitle string, articleBody string) (articleID string, e error)

	AddUserFeed(userID string, feedID string) error

	ListUserFeeds(userID string) ([]api.Feed, error)

	GetUserFeed(userID string, feedID string) (*api.Feed, error)

	ListUserArticles(userID string) ([]api.Article, error)

	ListUserFeedArticles(userID string, feedID string) ([]api.Article, error)

	Close()
}
