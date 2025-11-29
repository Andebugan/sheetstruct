package app

import "github.com/andebugan/sheetstruct/internal/models"

// Interface for Component actions
type IComponentManager interface {
	// Finds component by ID,
	// if sheet doesn't exist - returns nil and error
	Get(cid models.ID, sid models.ID, uid models.ID) (models.Component, error)

	// Returns collection of all components for specific sheet and user
	// if user does not exist, or has no access to sheet,
	// returns nil, error
	GetMany(sid models.ID, uid models.ID) ([]models.Component, error)

	// Creates new component for specific user
	// if creation is successfull - returns created component
	Create(component models.Component) (models.Component, error)

	// Deletes component by id,
	// if component does not exist, or user has no access to it,
	// returns nil, error
	Delete(cid models.ID, sid models.ID, uid models.ID) error

	// Finds component with matching id and updates it's values,
	// if user or component doesn' exit - returns error
	Update(sheet models.Sheet) error
}

// Component manager data
type ComponentManager struct{}

// Creates new component manager
func NewComponentManager() *ComponentManager {
	return &ComponentManager{}
}

// Finds component by ID,
// if sheet doesn't exist - returns nil and error
func (c *ComponentManager) Get(cid models.ID, sid models.ID, uid models.ID) (models.Component, error) {
	return models.Component{}, nil
}

// Returns collection of all components for specific sheet and user
// if user does not exist, or has no access to sheet,
// returns nil, error
func (c *ComponentManager) GetMany(sid models.ID, uid models.ID) ([]models.Component, error) {
	return []models.Component{}, nil
}

// Creates new component for specific user
// if creation is successfull - returns created component
func (c *ComponentManager) Create(component models.Component) (models.Component, error) {
	return models.Component{}, nil
}

// Deletes component by id,
// if component does not exist, or user has no access to it,
// returns nil, error
func (c* ComponentManager) Delete(cid models.ID, sid models.ID, uid models.ID) error {
	return nil
}

// Finds component with matching id and updates it's values,
// if user or component doesn' exit - returns error
func (c* ComponentManager) Update(sheet models.Sheet) error {
	return nil
}
