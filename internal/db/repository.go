package db

import (
	"github.com/if-ivan-else/tldrfeed/internal/types"
)

// Repository defines an interface with persistence layer for Users, Feeds, and Articles entities
type Repository interface {
	CreateUser(name string) (*types.User, error)

	ListUsers() ([]types.User, error)

	GetUserByID(userID string) (*types.User, error)

	CreateFeed(name string) (*types.Feed, error)

	ListFeeds() ([]types.Feed, error)

	GetFeedByID(feedID string) (*types.Feed, error)

	ListFeedArticles(feedID string) ([]types.Article, error)

	CreateFeedArticle(feedID string, articleTitle string, articleBody string) (articleID string, e error)

	AddUserFeed(userID string, feedID string) error

	ListUserFeeds(userID string) ([]types.Feed, error)

	ListUserArticles(userID string) ([]types.Article, error)

	ListUserFeedArticles(userID string, feedID string) ([]types.Article, error)
}

// NewRepository creates a new repository
func NewRepository() Repository {

	// TODO
	return nil
	// return &mock.Repository{}
}
