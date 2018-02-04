package mock

import (
	"github.com/google/uuid"
	"github.com/if-ivan-else/tldrfeed/api"
	"github.com/if-ivan-else/tldrfeed/internal/db"
)

// Repository is a mock repository implementation used in tests
type repository struct {
	users []api.User
	feeds []api.Feed

	userFeeds    map[string][]api.Feed
	feedArticles map[string][]api.Article
}

// NewRepository creates an instance of a mock repository for tests
func NewRepository() db.Repository {
	r := &repository{}
	r.userFeeds = make(map[string][]api.Feed)
	r.feedArticles = make(map[string][]api.Article)
	return r
}

func (r *repository) CreateUser(name string) (*api.User, error) {
	u := api.User{
		ID:   uuid.New().String(),
		Name: name,
	}

	r.users = append(r.users, u)
	r.userFeeds[u.ID] = []api.Feed{}
	return &u, nil
}

func (r *repository) ListUsers() ([]api.User, error) {
	return r.users, nil
}

func (r *repository) GetUser(userID string) (*api.User, error) {
	for _, u := range r.users {
		if u.ID == userID {
			return &u, nil
		}
	}
	return nil, db.ErrNoSuchUser
}

func (r *repository) CreateFeed(name string) (*api.Feed, error) {
	f := api.Feed{
		ID:   uuid.New().String(),
		Name: name,
	}
	r.feeds = append(r.feeds, f)
	r.feedArticles[f.ID] = []api.Article{}
	return &f, nil
}

func (r *repository) ListFeeds() ([]api.Feed, error) {
	return r.feeds, nil
}

func (r *repository) GetFeed(feedID string) (*api.Feed, error) {
	for _, f := range r.feeds {
		if f.ID == feedID {
			return &f, nil
		}
	}
	return nil, db.ErrNoSuchFeed
}

func (r *repository) ListFeedArticles(feedID string) ([]api.Article, error) {
	articles, ok := r.feedArticles[feedID]
	if !ok {
		return nil, db.ErrNoSuchFeed
	}
	return articles, nil
}

func (r *repository) CreateFeedArticle(feedID string, articleTitle string, articleBody string) (articleID string, e error) {

	a := api.Article{
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
	f, err := r.GetFeed(feedID)
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

func (r *repository) ListUserFeeds(userID string) ([]api.Feed, error) {
	feeds, ok := r.userFeeds[userID]
	if !ok {
		return nil, db.ErrNoSuchUser
	}
	return feeds, nil
}

func (r *repository) GetUserFeed(userID string, feedID string) (*api.Feed, error) {
	feeds, ok := r.userFeeds[userID]
	if !ok {
		return nil, db.ErrNoSuchUser
	}
	for _, f := range feeds {
		if f.ID == feedID {
			return &f, nil
		}
	}
	return nil, db.ErrNotSubscribed
}

func (r *repository) ListUserArticles(userID string) ([]api.Article, error) {
	feeds, err := r.ListUserFeeds(userID)
	if err != nil {
		return nil, err
	}

	userArticles := []api.Article{}
	for _, f := range feeds {
		articles, ok := r.feedArticles[f.ID]
		if ok {
			userArticles = append(userArticles, articles...)
		}
	}

	return userArticles, nil
}

func (r *repository) ListUserFeedArticles(userID string, feedID string) ([]api.Article, error) {
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

func (r *repository) Close() {
}
