package mongo

import (
	"github.com/google/uuid"
	"github.com/if-ivan-else/tldrfeed/internal/db"
	"github.com/if-ivan-else/tldrfeed/internal/types"
)

// import "github.com/globalsign/mgo"

// repository implements a MongoDB based repository for tldrfeed persistence of Users, Articles and Feeds
type repository struct {
}

// NewRepository creates an instance of a MongoDB repository for tests
func NewRepository(dbAddress string) db.Repository {
	r := &repository{}
	return r
}

func (r *repository) CreateUser(name string) (*types.User, error) {
	u := types.User{
		ID:   uuid.New().String(),
		Name: name,
	}

	// TODO
	return &u, nil
}

func (r *repository) ListUsers() ([]types.User, error) {
	// TODO
	return nil, db.ErrNotImplemented
	// return r.users, nil
}

func (r *repository) GetUser(userID string) (*types.User, error) {

	// TODO
	return nil, db.ErrNoSuchUser
}

func (r *repository) CreateFeed(name string) (*types.Feed, error) {
	f := types.Feed{
		ID:   uuid.New().String(),
		Name: name,
	}
	// TODO
	return &f, nil
}

func (r *repository) ListFeeds() ([]types.Feed, error) {
	// TODO
	return nil, nil
}

func (r *repository) GetFeed(feedID string) (*types.Feed, error) {
	// TODO
	return nil, db.ErrNoSuchFeed
}

func (r *repository) ListFeedArticles(feedID string) ([]types.Article, error) {
	// TODO
	return nil, db.ErrNotImplemented
}

func (r *repository) CreateFeedArticle(feedID string, articleTitle string, articleBody string) (articleID string, e error) {

	a := types.Article{
		ID:    uuid.New().String(),
		Title: articleTitle,
		Body:  articleBody,
	}

	return a.ID, nil
}

func (r *repository) AddUserFeed(userID string, feedID string) error {
	// f, err := r.GetFeed(feedID)
	// if err != nil {
	// 	return err
	// }

	// TODO
	return db.ErrNotImplemented
}

func (r *repository) ListUserFeeds(userID string) ([]types.Feed, error) {
	// TODO
	return nil, db.ErrNotImplemented
}

func (r *repository) GetUserFeed(userID string, feedID string) (*types.Feed, error) {
	return nil, db.ErrNotImplemented
}

func (r *repository) ListUserArticles(userID string) ([]types.Article, error) {
	return nil, db.ErrNotImplemented
}

func (r *repository) ListUserFeedArticles(userID string, feedID string) ([]types.Article, error) {
	return nil, db.ErrNotImplemented

}
