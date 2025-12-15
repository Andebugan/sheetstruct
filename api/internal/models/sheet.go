package models

import "time"

// Sheet information
type Sheet struct {
	Id            SheetID // Sheet Id
	UId           UserID // User Id
	Name          string
	Description   string
	Template      bool
	LastWriteTime time.Time   // Updates each time sheet is modified, used for sorting
	Components    []ComponentID // Array of component ID's connected to this sheet
}
