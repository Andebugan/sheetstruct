package dtos

import "github.com/andebugan/sheetstruct/internal/models"

// Database compatible DTO for user info
type User struct {
	Id           uint64 `gorm:"primaryKey" json:"id"`
	Name         string `gorm:"size:100;not null" json:"name"`
	Password     string `gorm:"size:100;not null" json:"password"`
	Email        string `gorm:"size:100;uniqueIndex;not null" json:"email"`
}

// Specified the table name for User DTO model
func (User) TableName() string {
	return "users"
}

// Function to convert model user to dto
func (u *User) ToDTO(user models.User) {
	u.Id = user.Id.Value
	u.Name = user.Name
	u.Password = user.Password
	u.Email = user.Email
}

// Function to convert model user to dto
func (u *User) FromDTO() models.User {
	return models.User{
		Id: models.UserID{ Value: u.Id },
		Name: u.Name,
		Password: u.Password,
		Email: u.Email,
	}
}
