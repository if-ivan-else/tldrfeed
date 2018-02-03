package types

// Article describes an Article posted to a Feed in the tldrfeed service
type Article struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}
