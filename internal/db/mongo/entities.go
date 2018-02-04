package mongo

import (
	"time"

	"github.com/if-ivan-else/tldrfeed/api"
)

// User is a Mongo document to store user records
type User struct {
	ID   string `bson:"_id"`
	Name string `bson:"name"`
}

func (u *User) toAPI() *api.User {
	return &api.User{
		ID:   u.ID,
		Name: u.Name,
	}
}

// UserList is a list of User documents
type UserList []User

func (l UserList) toAPI() []api.User {
	res := []api.User{}
	for _, u := range l {
		res = append(res, *u.toAPI())
	}
	return res
}

// Feed is a Mongo document to store feed records
type Feed struct {
	ID    string   `bson:"_id"`
	Name  string   `bson:"title"`
	Users []string `bson:"users"`
}

func (f *Feed) toAPI() *api.Feed {
	return &api.Feed{
		ID:   f.ID,
		Name: f.Name,
	}
}

// FeedList is a list of Feed documents
type FeedList []Feed

func (l FeedList) toAPI() []api.Feed {
	res := []api.Feed{}
	for _, f := range l {
		res = append(res, *f.toAPI())
	}
	return res
}

// Article is a Mongo document to store article records
type Article struct {
	ID            string    `bson:"_id"`
	FeedID        string    `bson:"feed_id"`
	Title         string    `bson:"title"`
	Body          string    `bson:"body"`
	PublishedTime time.Time `bson:"published_at"`
}

func (a *Article) toAPI() *api.Article {
	return &api.Article{
		ID:            a.ID,
		Title:         a.Title,
		Body:          a.Body,
		PublishedTime: a.PublishedTime,
	}
}

// ArticleList is a list of Article documents
type ArticleList []Article

func (l ArticleList) toAPI() []api.Article {
	res := []api.Article{}
	for _, a := range l {
		res = append(res, *a.toAPI())
	}
	return res
}
