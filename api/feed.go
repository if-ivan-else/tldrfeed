package api

// Feed defines a Feed in the tldrfeed service
type Feed struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CreateFeedRequest represents a request to create a new User
type CreateFeedRequest struct {
	Name string `json:"name" valid:"required~Feed name cannot be blank"`
}

// AddUserFeedRequest represents a request to subscribe a User to an existing Feed
type AddUserFeedRequest struct {
	FeedID string `json:"feed_id" valid:"required~Feed ID cannot be blank"`
}
