package managers

import (
	"context"

	"github.com/andebugan/sheetstruct/internal/app"
	"github.com/andebugan/sheetstruct/internal/app/dtos"
	"github.com/andebugan/sheetstruct/internal/app/repos"
	"github.com/andebugan/sheetstruct/internal/models"
)

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
}

// User manager data
type UserManager struct {
	ctx      context.Context
	userRepo repos.IUserRepository
}

// Initializes new User Manager
func NewUserManager(ctx context.Context, userRepo repos.IUserRepository) *UserManager {
	return &UserManager{
		ctx:      ctx,
		userRepo: userRepo,
	}
}

// Finds user by ID,
// if user doesn't exist - returns nil and error
func (u *UserManager) Get(uid models.UserID) (models.User, error) {
	user, err := u.userRepo.FindByID(u.ctx, uid.Value)

	if err != nil {
		return models.User{}, err
	}

	return user.FromDTO(), nil
}

// Returns collection of all users
func (u *UserManager) GetMany() ([]models.User, error) {
	userDtos, err := u.userRepo.GetAll(u.ctx, -1, -1)

	if err != nil {
		return []models.User{}, err
	}

	users := make([]models.User, len(userDtos))
	for i, user := range userDtos {
		users[i] = user.FromDTO()
	}

	return users, nil
}

// Registers new user
// if creation is successfull - returns created user
func (u *UserManager) Create(userData models.NewUserData) (models.User, error) {
	err := u.userRepo.CheckDuplicateName(u.ctx, userData.Name)
	if err != nil {
		return models.User{}, err
	}

	err = u.userRepo.CheckDuplicateEmail(u.ctx, userData.Name)
	if err != nil {
		return models.User{}, err
	}

	userDto := dtos.User{}
	userDto.ToDTO(models.User{
		Id: models.UserID{
			Value: 0,
		},
		Name:     userData.Name,
		Email:    userData.Email,
		Password: userData.Password,
	})

	user, err := u.userRepo.Create(u.ctx, &userDto)
	if err != nil {
		return models.User{}, err
	}

	return user.FromDTO(), nil
}

// Deletes user by id,
// if doesn't exist - returns error
func (u *UserManager) Delete(uid models.UserID) error {
	return u.userRepo.Delete(u.ctx, uid.Value)
}

// Finds user with matching id and updates it's values,
// if user doesn' exit - returns error
func (u *UserManager) Update(user models.User) (models.User, error) {
	userDto := dtos.User{}
	userDto.ToDTO(user)

	updatedDto, err := u.userRepo.Update(u.ctx, &userDto)

	if err != nil {
		return models.User{}, err
	}

	return updatedDto.FromDTO(), nil
}

// Returns user with matching credentials,
// if user doesn't exist or password is incorrect,
// returns nil and error
func (u *UserManager) Login(login string, password string) (models.User, error) {

	// Try find by name
	user, err := u.userRepo.FindByName(u.ctx, login)
	if err == nil && user.Password == password {
		return user.FromDTO(), nil
	}

	// Try find by email
	user, err = u.userRepo.FindByEmail(u.ctx, login)
	if err == nil && user.Password == password {
		return user.FromDTO(), nil
	}

	return models.User{}, app.ErrUserAuthFailed
}
