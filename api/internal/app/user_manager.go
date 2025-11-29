package app

import "github.com/andebugan/sheetstruct/internal/models"

// Interface for User actions
type IUserManager interface {
	// Finds user by ID,
	// if user doesn't exist - returns nil and error
	Get(id models.ID) (models.User, error)

	// Returns collection of all users
	GetMany() ([]models.User, error)

	// Registers new user
	// if creation is successfull - returns created user
	Create(user models.User) (models.User, error)

	// Deletes user by id,
	// if doesn't exist - returns error
	Delete(id models.ID) error

	// Finds user with matching id and updates it's values,
	// if user doesn' exit - returns error
	Update(user models.ID) error

	// Returns user with matching credentials,
	// if user doesn't exist or password is incorrect,
	// returns nil and error
	// TODO: not secure to pass pure strings, better replace with salted versions
	Auth(name string, password string) (models.User, error)
}

// User manager data
type UserManager struct {}


// Initializes new User Manager
func NewUserManager() *UserManager {
	return &UserManager{}
}

// Finds user by ID,
// if user doesn't exist - returns nil and error
func (u *UserManager) Get(id models.ID) (models.User, error) {
	return models.User{}, nil
}

// Returns collection of all users
func (u *UserManager) GetMany() ([]models.User, error) {
	return []models.User{}, nil
}

// Registers new user
// if creation is successfull - returns created user
func (u *UserManager) Create(user models.User) (models.User, error) {
	return models.User{}, nil
}

// Deletes user by id,
// if doesn't exist - returns error
func (u *UserManager) Delete(id models.ID) error {
	return nil
}

// Finds user with matching id and updates it's values,
// if user doesn' exit - returns error
func (u *UserManager) Update(user models.ID) error {
	return nil
}

// Returns user with matching credentials,
// if user doesn't exist or password is incorrect,
// returns nil and error
// TODO: not secure to pass pure strings, better replace with salted versions
func (u *UserManager) Auth(name string, password string) (models.User, error) {
	return models.User{}, nil
}

