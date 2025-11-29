package app

import "github.com/andebugan/sheetstruct/internal/models"

// Interface for Sheet actions
type ISheetManager interface {
	// Finds Sheet by ID,
	// if sheet doesn't exist - returns nil and error
	Get(sid models.ID, uid models.ID) (models.Sheet, error)

	// Returns collection of all sheets, accessible to user,
	// if user does not exist, or has no access to sheet,
	// returns nil, error
	GetMany(uid models.ID) ([]models.Sheet, error)

	// Creates new sheet
	// if creation is successfull - returns created sheet
	Create(sheet models.Sheet) (models.Sheet, error)

	// Deletes sheet, accessible to user, by id,
	// if user does not exist, or has no access to sheet,
	// returns nil, error
	Delete(sid models.ID, uid models.ID) error

	// Finds sheet with matching id and updates it's values,
	// if user or sheet doesn' exit - returns error
	Update(sheet models.Sheet) error
}

// Sheet manager data
type SheetManager struct{}

// Creates new sheet manager
func NewSheetManager() *SheetManager {
	return &SheetManager{}
}

// Finds Sheet by ID,
// if sheet doesn't exist - returns nil and error
func (s* SheetManager) Get(sid models.ID, uid models.ID) (models.Sheet, error) {
	return models.Sheet{}, nil
}

// Returns collection of all sheets, accessible to user,
// if user does not exist, or has no access to sheet,
// returns nil, error
func (s* SheetManager) GetMany(uid models.ID) ([]models.Sheet, error) {
	return []models.Sheet{}, nil
}

// Creates new sheet
// if creation is successfull - returns created sheet
func (s* SheetManager) Create(sheet models.Sheet) (models.Sheet, error) {
	return models.Sheet{}, nil
}

// Deletes sheet, accessible to user, by id,
// if user does not exist, or has no access to sheet,
// returns nil, error
func (s* SheetManager) Delete(sid models.ID, uid models.ID) error {
	return nil
}

// Finds sheet with matching id and updates it's values,
// if user or sheet doesn' exit - returns error
func (s* SheetManager) Update(sheet models.Sheet) error {
	return nil
}

