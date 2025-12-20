package dtos

import "github.com/andebugan/sheetstruct/internal/models"

// Database compatible DTO for user info
type User struct {
	Id           uint64 `gorm:"primaryKey" json:"id"`
	Name         string `json:"name"`
	Password     string `json:"password"`
	Email        string `gorm:"uniqueIndex;not null" json:"email"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) ToDTO(user models.User) {
	u.Id = user.Id.Value
	u.Name = user.Name
	u.Password = user.Password
	u.Email = user.Email
}

func (u *User) FromDTO() models.User {
	return models.User{
		Id: models.UserID{ Value: u.Id },
		Name: u.Name,
		Password: u.Password,
		Email: u.Email,
	}
}
