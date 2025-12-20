package dtos

import "github.com/andebugan/sheetstruct/internal/models"

// Acts as building block of each sheet
type Component struct {
	Id          uint64 `gorm:"primaryKey" json:"id"`
	UId         uint64 `gorm:"foreignKey" json:"uid"`
	SId         uint64 `gorm:"foreignKey" json:"sid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Template    bool   `gorm:"not null;default:false" json:"template"`
	VarName     string `json:"varname"`
	VarType     int    `json:"vartype"`
	VarValue    []byte `json:"varvalue"`
	Style       string `json:"style"`
}

func (Component) TableName() string {
	return "components"
}

func (c *Component) ToDTO(component models.Component) {
	c.Id = component.Id.Value
	c.SId = component.SId.Value
	c.UId = component.UId.Value
	c.Name = component.Name
	c.Description = component.Description
	c.Template = component.Template
	c.VarName = component.VarName
	c.VarType = int(component.VarType)
	c.VarValue = component.VarValue
	c.Style = string(component.Style)
}

func (c *Component) FromDTO() models.Component {
	return models.Component{
		Id: models.ComponentID{ Value: c.Id },
		SId: models.SheetID{ Value: c.SId },
		UId: models.UserID{ Value: c.UId },
		Name: c.Name,
		Description: c.Description,
		Template: c.Template,
		VarName: c.VarName,
		VarType: models.VariableType(c.VarType),
		VarValue: c.VarValue,
		Style: models.Style(c.Style),
	}
}
