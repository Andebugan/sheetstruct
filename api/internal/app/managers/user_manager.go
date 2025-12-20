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
	Get(uid models.UserID) (*models.User, error)

	// Returns collection of all users
	GetMany() ([]models.User, error)

	// Registers new user,
	// if creation is successfull - returns created user
	Create(userData models.NewUserData) (*models.User, error)

	// Deletes user by id,
	// if doesn't exist - returns error
	Delete(uid models.UserID) error

	// Finds user with matching id and updates it's values,
	// if user doesn' exit - returns error
	Update(user models.User) (*models.User, error)

	// Returns user with matching credentials,
	// username and email can be both used as login,
	// if user doesn't exist or password is incorrect,
	// returns nil and error
	Login(login string, password string) (*models.User, error)
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
func (m *UserManager) Get(uid models.UserID) (*models.User, error) {
	userDto, err := m.userRepo.FindByID(m.ctx, uid.Value)

	if err != nil {
		return nil, err
	}

	user := userDto.FromDTO()

	return &user, nil
}

// Returns collection of all users
func (m *UserManager) GetMany() ([]models.User, error) {
	userDtos, err := m.userRepo.GetAll(m.ctx, -1, -1)

	if err != nil {
		return nil, err
	}

	users := make([]models.User, len(userDtos))
	for i, user := range userDtos {
		users[i] = user.FromDTO()
	}

	return users, nil
}

// Registers new user
// if creation is successfull - returns created user
func (m *UserManager) Create(userData models.NewUserData) (*models.User, error) {
	err := m.userRepo.CheckDuplicateName(m.ctx, userData.Name)
	if err != nil {
		return nil, err
	}

	err = m.userRepo.CheckDuplicateEmail(m.ctx, userData.Email)
	if err != nil {
		return nil, err
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

	newUserDto, err := m.userRepo.Create(m.ctx, &userDto)
	if err != nil {
		return nil, err
	}

	user := newUserDto.FromDTO()

	return &user, nil
}

// Deletes user by id
func (m *UserManager) Delete(uid models.UserID) error {
	return m.userRepo.Delete(m.ctx, uid.Value)
}

// Finds user with matching id and updates it's values,
// if user doesn' exit - returns error
func (m *UserManager) Update(user models.User) (*models.User, error) {
	userDto := dtos.User{}
	userDto.ToDTO(user)

	updatedDto, err := m.userRepo.Update(m.ctx, &userDto)

	if err != nil {
		return nil, err
	}

	updatedUser := updatedDto.FromDTO()

	return &updatedUser, nil
}

// Returns user with matching credentials,
// if user doesn't exist or password is incorrect,
// returns nil and error
func (m *UserManager) Login(login string, password string) (*models.User, error) {
	// Try find by name
	userDto, err := m.userRepo.FindByName(m.ctx, login)
	if err == nil && userDto.Password == password {
		user := userDto.FromDTO()
		return &user, nil
	}

	// Try find by email
	userDto, err = m.userRepo.FindByEmail(m.ctx, login)
	if err == nil && userDto.Password == password {
		user := userDto.FromDTO()
		return &user, nil
	}

	return nil, app.ErrUserAuthFailed
}
