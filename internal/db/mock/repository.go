package mock

import (
	"github.com/google/uuid"
	"github.com/if-ivan-else/tldrfeed/internal/db"
	"github.com/if-ivan-else/tldrfeed/internal/types"
)

// Repository is a mock repository implementation
type repository struct {
	users []types.User
	feeds []types.Feed

	userFeeds    map[string][]types.Feed
	feedArticles map[string][]types.Article
}

func NewRepository() db.Repository {
	r := &repository{}
	r.userFeeds = make(map[string][]types.Feed)
	r.feedArticles = make(map[string][]types.Article)
	return r
}

func (r *repository) CreateUser(name string) (*types.User, error) {
	u := types.User{
		ID:   uuid.New().String(),
		Name: name,
	}

	r.users = append(r.users, u)
	r.userFeeds[u.ID] = []types.Feed{}
	return &u, nil
}

func (r *repository) ListUsers() ([]types.User, error) {
	return r.users, nil
}

func (r *repository) GetUserByID(userID string) (*types.User, error) {
	for _, u := range r.users {
		if u.ID == userID {
			return &u, nil
		}
	}
	return nil, db.ErrNoSuchUser
}

func (r *repository) CreateFeed(name string) (*types.Feed, error) {
	f := types.Feed{
		ID:   uuid.New().String(),
		Name: name,
	}
	r.feeds = append(r.feeds, f)
	return &f, nil
}

func (r *repository) ListFeeds() ([]types.Feed, error) {
	return r.feeds, nil
}

func (r *repository) GetFeedByID(feedID string) (*types.Feed, error) {
	for _, f := range r.feeds {
		if f.ID == feedID {
			return &f, nil
		}
	}
	return nil, db.ErrNoSuchFeed
}

func (r *repository) ListFeedArticles(feedID string) ([]types.Article, error) {
	articles, ok := r.feedArticles[feedID]
	if !ok {
		return nil, db.ErrNoSuchFeed
	}
	return articles, nil
}

func (r *repository) CreateFeedArticle(feedID string, articleTitle string, articleBody string) (articleID string, e error) {

	a := types.Article{
		ID:    uuid.New().String(),
		Title: articleTitle,
		Body:  articleBody,
	}

	articles, ok := r.feedArticles[feedID]
	if !ok {
		return "", db.ErrNoSuchFeed
	}

	r.feedArticles[feedID] = append(articles, a)
	return a.ID, nil
}

func (r *repository) AddUserFeed(userID string, feedID string) error {
	f, err := r.GetFeedByID(feedID)
	if err != nil {
		return err
	}

	feeds, ok := r.userFeeds[userID]
	if !ok {
		return db.ErrNoSuchUser
	}

	r.userFeeds[userID] = append(feeds, *f)

	return nil
}

func (r *repository) ListUserFeeds(userID string) ([]types.Feed, error) {
	feeds, ok := r.userFeeds[userID]
	if !ok {
		return nil, db.ErrNoSuchUser
	}
	return feeds, nil
}

func (r *repository) ListUserArticles(userID string) ([]types.Article, error) {
	feeds, err := r.ListUserFeeds(userID)
	if err != nil {
		return nil, err
	}

	userArticles := []types.Article{}
	for _, f := range feeds {
		articles, ok := r.feedArticles[f.ID]
		if ok {
			userArticles = append(userArticles, articles...)
		}
	}

	return userArticles, nil
}

func (r *repository) ListUserFeedArticles(userID string, feedID string) ([]types.Article, error) {
	feeds, err := r.ListUserFeeds(userID)
	if err != nil {
		return nil, err
	}

	for _, f := range feeds {
		if f.ID == feedID {
			return r.ListFeedArticles(feedID)
		}
	}
	return nil, db.ErrNotSubscribed

}
