package models

import "time"

// Sheet information
type Sheet struct {
	Id            SheetID // Sheet Id
	UId           UserID // User Id
	Name          string
	Description   string
	Template      bool
	LastEditTime time.Time   // Updates each time sheet is modified, used for sorting
}

func NewDefaultSheet(uid UserID) Sheet {
	return Sheet{
		Id: SheetID{ Value: 0 },
		UId: uid,
		Name: "sheet",
		Description: "",
		Template: false,
		LastEditTime: time.Now(),
	}
}
