package api

import "time"

// Article describes an Article posted to a Feed in the tldrfeed service
type Article struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Body          string    `json:"body"`
	PublishedTime time.Time `json:"published_at"`
}
