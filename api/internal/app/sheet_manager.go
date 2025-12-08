package app

import (
	"errors"

	"github.com/andebugan/sheetstruct/internal/models"
)

var ErrSheetNotFound = errors.New("Unable to find requested user")
var ErrSheetAlredyExists = errors.New("User already exists")

// Interface for Sheet actions
type ISheetManager interface {
	// Finds Sheet by ID,
	// if sheet doesn't exist - returns nil and error
	Get(sid models.SheetID, uid models.UserID) (models.Sheet, error)

	// Returns collection of all sheets, accessible to user,
	// if user does not exist, or has no access to sheet,
	// returns nil, error
	GetMany(uid models.UserID) ([]models.Sheet, error)

	// Creates new sheet with default values
	// if creation is successfull - returns created sheet
	Create(uid models.UserID) (models.Sheet, error)

	// Deletes sheet, accessible to user, by id,
	// if user does not exist, or has no access to sheet,
	// returns nil, error
	Delete(sid models.SheetID, uid models.UserID) error

	// Finds sheet with matching id and updates it's values,
	// if user or sheet doesn' exit - returns error
	Update(sheet models.Sheet) (models.Sheet, error)
}

// Sheet manager data
type SheetManager struct{}

// Creates new sheet manager
func NewSheetManager() *SheetManager {
	return &SheetManager{}
}

// Finds Sheet by ID,
// if sheet doesn't exist - returns nil and error
func (s *SheetManager) Get(sid models.SheetID, uid models.UserID) (models.Sheet, error) {
	return models.Sheet{}, nil
}

// Returns collection of all sheets, accessible to user,
// if user does not exist, or has no access to sheet,
// returns nil, error
func (s *SheetManager) GetMany(uid models.UserID) ([]models.Sheet, error) {
	return []models.Sheet{}, nil
}

// Creates new sheet
// if creation is successfull - returns created sheet
func (s *SheetManager) Create(uid models.UserID) (models.Sheet, error) {
	return models.Sheet{}, nil
}

// Deletes sheet, accessible to user, by id,
// if user does not exist, or has no access to sheet,
// returns nil, error
func (s *SheetManager) Delete(sid models.SheetID, uid models.UserID) error {
	return nil
}

// Finds sheet with matching id and updates it's values,
// if user or sheet doesn' exit - returns error
func (s *SheetManager) Update(sheet models.Sheet) (models.Sheet, error) {
	return models.Sheet{}, nil
}
