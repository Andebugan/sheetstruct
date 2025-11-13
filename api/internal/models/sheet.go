package models

import "time"

// Sheet information
type Sheet struct {
	Id ID // Sheet Id
	UId ID // User Id
	Name string
	Description string
	Template bool
	LastWriteTime time.Time // Updates each time sheet is modified, used for sorting
	Components []Component // Array of components, connected to this sheet
}

