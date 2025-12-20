package managers

// TODO: implement user/sheet access management

import (
	"context"

	"github.com/andebugan/sheetstruct/internal/app/dtos"
	"github.com/andebugan/sheetstruct/internal/app/repos"
	"github.com/andebugan/sheetstruct/internal/models"
)

// Interface for Component actions
type IComponentManager interface {
	// Finds component by ID,
	// if sheet doesn't exist - returns nil and error
	Get(cid models.ComponentID, sid models.SheetID, uid models.UserID) (*models.Component, error)

	// Clones component by ID,
	// if sheet doesn't exist - returns nil and error
	Clone(cid models.ComponentID, sid models.SheetID, uid models.UserID) (*models.Component, error)

	// Returns collection of all components for specific sheet and user
	// if user does not exist, or has no access to sheet,
	// returns nil, error
	GetMany(sid models.SheetID, uid models.UserID) ([]models.Component, error)

	// Creates new component for specific user
	// if creation is successfull - returns created component
	Create(sid models.SheetID, uid models.UserID) (*models.Component, error)

	// Deletes component by id,
	// if component does not exist, or user has no access to it,
	// returns nil, error
	Delete(cid models.ComponentID, sid models.SheetID, uid models.UserID) error

	// Finds component with matching id and updates it's values,
	// if user or component doesn' exit - returns error
	Update(sheet models.Component, uid models.UserID) (*models.Component, error)
}

// Component manager data
type ComponentManager struct {
	ctx           context.Context
	userRepo      repos.IUserRepository
	sheetRepo     repos.ISheetRepository
	componentRepo repos.IComponentRepository
}

// Creates new component manager
func NewComponentManager(ctx context.Context, userRepo repos.IUserRepository, sheetRepo repos.ISheetRepository, componentRepo repos.IComponentRepository) *ComponentManager {
	return &ComponentManager{
		ctx:           ctx,
		userRepo:      userRepo,
		sheetRepo:     sheetRepo,
		componentRepo: componentRepo,
	}
}

// Finds component by ID,
// if sheet doesn't exist - returns nil and error
func (m *ComponentManager) Get(cid models.ComponentID, sid models.SheetID, uid models.UserID) (*models.Component, error) {
	dto, err := m.componentRepo.FindByID(m.ctx, cid.Value)

	if err != nil {
		return nil, err
	}

	component := dto.FromDTO()

	return &component, nil
}

// Clones component by ID,
// if sheet doesn't exist - returns nil and error
func (m *ComponentManager) Clone(cid models.ComponentID, sid models.SheetID, uid models.UserID) (*models.Component, error) {
	component, err := m.Get(cid, sid, uid)
	if err != nil {
		return nil, err
	}

	dto := dtos.Component{}
	dto.ToDTO(*component)

	cloneDto, err := m.componentRepo.Create(m.ctx, &dto)
	if err != nil {
		return nil, err
	}

	cloneComponent := cloneDto.FromDTO()

	return &cloneComponent, nil
}

// Returns collection of all components for specific sheet and user
// if user does not exist, or has no access to sheet,
// returns nil, error
func (m *ComponentManager) GetMany(sid models.SheetID, uid models.UserID) ([]models.Component, error) {
	dtos, err := m.componentRepo.FindForSheet(m.ctx, sid.Value)

	if err != nil {
		return nil, err
	}

	components:= make([]models.Component, len(dtos))
	for i, component:= range dtos {
		components[i] = component.FromDTO()
	}

	return components, nil
}

// Creates new component for specific user
// if creation is successfull - returns created component
func (m *ComponentManager) Create(sid models.SheetID, uid models.UserID) (*models.Component, error) {
	dto := dtos.Component{}
	dto.ToDTO(models.NewDefaultComponent(sid, uid))
	
	newDto, err := m.componentRepo.Create(m.ctx, &dto)
	if err != nil {
		return nil, err
	}

	newComponent := newDto.FromDTO()

	return &newComponent, nil
}

// Deletes component by id,
// if component does not exist, or user has no access to it,
// returns nil, error
func (m *ComponentManager) Delete(cid models.ComponentID, sid models.SheetID, uid models.UserID) error {
	return m.componentRepo.Delete(m.ctx, sid.Value)
}

// Finds component with matching id and updates it's values,
// if user or component doesn' exit - returns error
func (m *ComponentManager) Update(component models.Component, uid models.UserID) (*models.Component, error) {
	dto := dtos.Component{}
	dto.ToDTO(component)

	updatedDto, err := m.componentRepo.Update(m.ctx, &dto)

	if err != nil {
		return nil, err
	}

	updatedComponent := updatedDto.FromDTO()

	return &updatedComponent, nil
}
