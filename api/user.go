package api

// User describes a User in the system
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CreateUserRequest represents a request to create a new User
type CreateUserRequest struct {
	Name string `json:"name" valid:"required~User name cannot be blank,alphanum~User name should be alphanumeric"`
}
