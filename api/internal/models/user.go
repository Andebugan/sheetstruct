package models

// Model of user info
type User struct {
	Id       ID // User Id
	Name     string
	Password string
	Email    string
}

