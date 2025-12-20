package dtos

import (
	"time"

	"github.com/andebugan/sheetstruct/internal/models"
)

// Database compatible DTO for sheet info
type Sheet struct {
	Id           uint64    `gorm:"primaryKey" json:"id"`
	UId          uint64    `gorm:"foreignKey" json:"uid"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Template     bool      `gorm:"not null;default:false" json:"template"`
	LastEditTime time.Time `gorm:"not null" json:"lastedittime"`
}

func (Sheet) TableName() string {
	return "sheets"
}

func (s *Sheet) ToDTO(sheet models.Sheet) {
	s.Id = sheet.Id.Value
	s.UId = sheet.UId.Value
	s.Name = sheet.Name
	s.Description = sheet.Description
	s.Template = sheet.Template
	s.LastEditTime = sheet.LastEditTime
}

func (s *Sheet) FromDTO() models.Sheet {
	return models.Sheet{
		Id: models.SheetID{ Value: s.Id },
		UId: models.UserID{ Value: s.UId },
		Name: s.Name,
		Description: s.Description,
		Template: s.Template,
		LastEditTime: s.LastEditTime,
	}
}
