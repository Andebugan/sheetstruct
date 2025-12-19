package managers

// TODO: implement user/sheet access management

import (
	"context"
	"time"

	"github.com/andebugan/sheetstruct/internal/app/dtos"
	"github.com/andebugan/sheetstruct/internal/app/repos"
	"github.com/andebugan/sheetstruct/internal/models"
)

// Interface for Sheet actions
type ISheetManager interface {
	// Finds Sheet by ID,
	// if sheet doesn't exist - returns nil and error
	Get(sid models.SheetID, uid models.UserID) (*models.Sheet, error)

	// Returns collection of all sheets, accessible to user,
	// if user does not exist, or has no access to sheet,
	// returns nil, error
	GetMany(uid models.UserID) ([]models.Sheet, error)

	// Creates new sheet with default values
	// if creation is successfull - returns created sheet
	Create(uid models.UserID) (*models.Sheet, error)

	// Deletes sheet, accessible to user, by id,
	// if user does not exist, or has no access to sheet,
	// returns nil, error
	Delete(sid models.SheetID, uid models.UserID) error

	// Finds sheet with matching id and updates it's values,
	// if user or sheet doesn' exit - returns error
	Update(sheet models.Sheet, uid models.UserID) (*models.Sheet, error)

	// Finds sheet with matching id and creates it's deep clone
	// if user or sheet doesn' exit - returns error
	Clone(sid models.SheetID, uid models.UserID) (*models.Sheet, error)
}

// Sheet manager data
type SheetManager struct {
	ctx       context.Context
	userRepo  repos.IUserRepository
	sheetRepo repos.ISheetRepository
}

// Creates new sheet manager
func NewSheetManager(ctx context.Context, userRepo repos.IUserRepository, sheetRepo repos.ISheetRepository) *SheetManager {
	return &SheetManager{
		ctx:       ctx,
		userRepo:  userRepo,
		sheetRepo: sheetRepo,
	}
}

// Finds Sheet by ID,
// if sheet doesn't exist - returns nil and error
func (m *SheetManager) Get(sid models.SheetID, uid models.UserID) (*models.Sheet, error) {
	dto, err := m.sheetRepo.FindByID(m.ctx, sid.Value)

	if err != nil {
		return nil, err
	}

	sheet := dto.FromDTO()

	return &sheet, nil
}

// Returns filtered collection of all sheets, accessible to user,
// if user does not exist, or has no access to sheet,
// returns nil, error
func (m *SheetManager) GetMany(uid models.UserID) ([]models.Sheet, error) {
	dtos, err := m.sheetRepo.FindByUser(m.ctx, uid.Value)

	if err != nil {
		return nil, err
	}

	sheets := make([]models.Sheet, len(dtos))
	for i, sheet := range dtos {
		sheets[i] = sheet.FromDTO()
	}

	return sheets, nil
}

// Creates new sheet
// if creation is successfull - returns created sheet
func (m *SheetManager) Create(uid models.UserID) (*models.Sheet, error) {
	dto := dtos.Sheet{}
	dto.ToDTO(models.NewDefaultSheet(uid))
	
	newDto, err := m.sheetRepo.Create(m.ctx, &dto)
	if err != nil {
		return nil, err
	}

	newSheet := newDto.FromDTO()

	return &newSheet, nil
}

// Deletes sheet, accessible to user, by id,
// if user does not exist, or has no access to sheet,
// returns nil, error
func (m *SheetManager) Delete(sid models.SheetID, uid models.UserID) error {
	return m.sheetRepo.Delete(m.ctx, sid.Value)
}

// Finds sheet with matching id and updates it's values,
// if user or sheet doesn' exit - returns error
func (m *SheetManager) Update(sheet models.Sheet, uid models.UserID) (*models.Sheet, error) {
	dto := dtos.Sheet{}
	dto.ToDTO(sheet)

	updatedDto, err := m.sheetRepo.Update(m.ctx, &dto)

	if err != nil {
		return nil, err
	}

	updatedSheet := updatedDto.FromDTO()

	return &updatedSheet, nil
}

// Finds sheet with matching id and creates it's deep clone
// if user or sheet doesn' exit - returns error
func (m *SheetManager) Clone(sid models.SheetID, uid models.UserID) (*models.Sheet, error) {
	sheet, err := m.Get(sid, uid)
	if err != nil {
		return nil, err
	}

	dto := dtos.Sheet{}
	dto.ToDTO(*sheet)

	// Reset ID and update UId to current user for clone
	// Update name to indicate it's a clone
	dto.Id = 0
	dto.UId = uid.Value
	if dto.Name != "" {
		dto.Name = dto.Name + " (копия)"
	}
	dto.LastEditTime = time.Now()

	cloneDto, err := m.sheetRepo.Create(m.ctx, &dto)
	if err != nil {
		return nil, err
	}

	cloneSheet := cloneDto.FromDTO()

	return &cloneSheet, nil
}
