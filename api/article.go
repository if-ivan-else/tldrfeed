package api

import "time"

// Article describes an Article posted to a Feed in the tldrfeed service
type Article struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Body          string    `json:"body"`
	PublishedTime time.Time `json:"published_at"`
}

// CreateArticleRequest defines a request to add an Article to a Feed
type CreateArticleRequest struct {
	Title string `json:"title" valid:"required~Article title cannot be blank"`
	Body  string `json:"body" valid:"required~Article title cannot be blank"`
}

// CreateArticleResponse defines a response to send for adding an Article to a Feed
type CreateArticleResponse struct {
	ID string `json:"id"`
}
