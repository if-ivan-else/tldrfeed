// Package api implements REST Client API for the tldrfeed service
package api

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
)

// Client implements a REST Client for programmatic interaction with tldrfeed service
type Client struct {
	sling *sling.Sling
}

// NewClient returns a new API client for tldrfeed
func NewClient(url string) *Client {

	httpClient := http.DefaultClient
	baseURL := fmt.Sprintf("%s%s/", url, APIVersion)
	return &Client{
		sling: sling.New().Client(httpClient).Base(baseURL),
	}
}

// CreateUser creates a new User
func (c *Client) CreateUser(name string) (*User, error) {

	createUser := &CreateUserRequest{
		Name: name,
	}

	var u User
	_, err := c.sling.Post("users").BodyJSON(createUser).ReceiveSuccess(&u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// CreateFeed creates a new Feed
func (c *Client) CreateFeed(name string) (*Feed, error) {

	createFeed := &CreateFeedRequest{
		Name: name,
	}
	var f Feed
	_, err := c.sling.Post("feeds").BodyJSON(createFeed).ReceiveSuccess(&f)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

// CreateArticle creates a new Article
func (c *Client) CreateArticle(feedID string, title string, body string) (*Article, error) {

	createArticle := &CreateArticleRequest{
		Title: title,
		Body:  body,
	}
	var f Article
	_, err := c.sling.Post(fmt.Sprintf("feeds/%s", feedID)).BodyJSON(createArticle).ReceiveSuccess(&f)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

// ListUsers lists all Users
func (c *Client) ListUsers() ([]User, error) {
	users := []User{}
	_, err := c.sling.Get("users").ReceiveSuccess(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// ListFeeds lists all Feeds
func (c *Client) ListFeeds() ([]Feed, error) {
	feeds := []Feed{}
	_, err := c.sling.Get("feeds").ReceiveSuccess(&feeds)
	if err != nil {
		return nil, err
	}
	return feeds, nil
}

// ListArticles lists all Articles
func (c *Client) ListArticles(feedID string) ([]Article, error) {
	articles := []Article{}
	_, err := c.sling.Get(fmt.Sprintf("feeds/%s/articles", feedID)).ReceiveSuccess(&articles)
	if err != nil {
		return nil, err
	}
	return articles, nil
}

// ListUserArticles lists Articles from all or one channel for a User
func (c *Client) ListUserArticles(userID string, feedID string) ([]Article, error) {
	articles := []Article{}

	var url string
	if feedID == "" {
		url = fmt.Sprintf("users/%s/articles", userID)
	} else {
		url = fmt.Sprintf("users/%s/feeds/%s/articles", userID, feedID)
	}
	_, err := c.sling.Get(url).ReceiveSuccess(&articles)
	if err != nil {
		return nil, err
	}
	return articles, nil
}
