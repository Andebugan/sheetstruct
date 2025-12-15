package app

import (
	"errors"

	"github.com/andebugan/sheetstruct/internal/models"
)

var ErrComponentNotFound = errors.New("Unable to find requested user")
var ErrComponentAlredyExists = errors.New("User already exists")

// Interface for Component actions
type IComponentManager interface {
	// Finds component by ID,
	// if sheet doesn't exist - returns nil and error
	Get(cid models.ComponentID, sid models.SheetID, uid models.UserID) (models.Component, error)

	// Clones component by ID,
	// if sheet doesn't exist - returns nil and error
	Clone(cid models.ComponentID, sid models.SheetID, uid models.UserID) (models.Component, error)

	// Returns collection of all components for specific sheet and user
	// if user does not exist, or has no access to sheet,
	// returns nil, error
	GetMany(sid models.SheetID, uid models.UserID) ([]models.Component, error)

	// Creates new component for specific user
	// if creation is successfull - returns created component
	Create(sid models.SheetID, uid models.UserID) (models.Component, error)

	// Deletes component by id,
	// if component does not exist, or user has no access to it,
	// returns nil, error
	Delete(cid models.ComponentID, sid models.SheetID, uid models.UserID) error

	// Finds component with matching id and updates it's values,
	// if user or component doesn' exit - returns error
	Update(sheet models.Component, uid models.UserID) (models.Component, error)
}

// Component manager data
type ComponentManager struct{}

// Creates new component manager
func NewComponentManager() *ComponentManager {
	return &ComponentManager{}
}

// Finds component by ID,
// if sheet doesn't exist - returns nil and error
func (c *ComponentManager) Get(cid models.ComponentID, sid models.SheetID, uid models.UserID) (models.Component, error) {
	return models.Component{}, nil
}

// Clones component by ID,
// if sheet doesn't exist - returns nil and error
func (c *ComponentManager) Clone(cid models.ComponentID, sid models.SheetID, uid models.UserID) (models.Component, error) {
	return models.Component{}, nil
}

// Returns collection of all components for specific sheet and user
// if user does not exist, or has no access to sheet,
// returns nil, error
func (c *ComponentManager) GetMany(sid models.SheetID, uid models.UserID) ([]models.Component, error) {
	return []models.Component{}, nil
}

// Creates new component for specific user
// if creation is successfull - returns created component
func (c *ComponentManager) Create(sid models.SheetID, uid models.UserID) (models.Component, error) {
	return models.Component{}, nil
}

// Deletes component by id,
// if component does not exist, or user has no access to it,
// returns nil, error
func (c *ComponentManager) Delete(cid models.ComponentID, sid models.SheetID, uid models.UserID) error {
	return nil
}

// Finds component with matching id and updates it's values,
// if user or component doesn' exit - returns error
func (c *ComponentManager) Update(sheet models.Component, uid models.UserID) (models.Component, error) {
	return models.Component{}, nil
}
