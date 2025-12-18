package models

// Model of info for initializing new user
type NewUserData struct {
	Name     string
	Password string
	Email    string
}

// Model of user info
type User struct {
	Id           UserID `gorm:"primarykey" json:"id"`
	Name         string `gorm:"size:100;not null" json:"name"`
	Password     string `gorm:"size:100;not null" json:"password"`
	Email        string `gorm:"size:100;uniqueindex;not null" json:"email"`
}
