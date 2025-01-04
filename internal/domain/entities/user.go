package entities

// User is an entity that represents a user
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}
