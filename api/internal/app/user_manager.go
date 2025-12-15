package app

import (
	"errors"

	"github.com/andebugan/sheetstruct/internal/models"
)

var ErrUserAuthFailed = errors.New("Unable to find user with matching credentials")
var ErrUserNotFound = errors.New("Unable to find requested user")
var ErrUserAlredyExists = errors.New("User already exists")

// Interface for User actions
type IUserManager interface {
	// Finds user by ID,
	// if user doesn't exist - returns nil and error
	Get(uid models.UserID) (models.User, error)

	// Returns collection of all users
	GetMany() ([]models.User, error)

	// Registers new user,
	// if creation is successfull - returns created user
	Create(userData models.NewUserData) (models.User, error)

	// Deletes user by id,
	// if doesn't exist - returns error
	Delete(uid models.UserID) error

	// Finds user with matching id and updates it's values,
	// if user doesn' exit - returns error
	Update(user models.User) (models.User, error)

	// Returns user with matching credentials,
	// username and email can be both used as login,
	// if user doesn't exist or password is incorrect,
	// returns nil and error
	Login(login string, password string) (models.User, error)

	// Logs user out of the system by clearing refresh token
	Logout(uid models.UserID) (error)
}

// User manager data
type UserManager struct{}

// Initializes new User Manager
func NewUserManager() *UserManager {
	return &UserManager{}
}

// Finds user by ID,
// if user doesn't exist - returns nil and error
func (u *UserManager) Get(uid models.UserID) (models.User, error) {
	return models.User{}, nil
}

// Returns collection of all users
func (u *UserManager) GetMany() ([]models.User, error) {
	return []models.User{}, nil
}

// Registers new user
// if creation is successfull - returns created user
func (u *UserManager) Create(userData models.NewUserData) (models.User, error) {
	return models.User{}, nil
}

// Deletes user by id,
// if doesn't exist - returns error
func (u *UserManager) Delete(uid models.UserID) error {
	return nil
}

// Finds user with matching id and updates it's values,
// if user doesn' exit - returns error
func (u *UserManager) Update(user models.User) (models.User, error) {
	return models.User{}, nil
}

// Returns user with matching credentials,
// if user doesn't exist or password is incorrect,
// returns nil and error
func (u *UserManager) Login(login string, password string) (models.User, error) {
	return models.User{}, nil
}

// Logs user out of the system by clearing refresh token
func (u *UserManager) Logout(uid models.UserID) (error) {
	//user.RefreshToken = ""
	return nil
}
