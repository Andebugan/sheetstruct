package models

// Model of info for initializing new user
type NewUserData struct {
	Name     string
	Password string
	Email    string
}

// Model of user info
type User struct {
	Id       UserID // User Id
	Name     string
	Password string
	Email    string
}
