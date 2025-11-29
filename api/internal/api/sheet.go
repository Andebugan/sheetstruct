package api

import (
	"github.com/andebugan/sheetstruct/internal/app"
	"github.com/gin-gonic/gin"
)

// Creates handler func for getting multiple sheets
func NewSheetsGetHandler(sheetManager app.ISheetManager) func(c *gin.Context) {
	return func(c *gin.Context) {
	}
}

// Creates handler func for getting single sheet
func NewSheetGetHandler(sheetManager app.ISheetManager) func(c *gin.Context) {
	return func(c *gin.Context) {
	}
}

// Creates handler func for creating new sheet
func NewSheetCreateHandler(sheetManager app.ISheetManager) func(c *gin.Context) {
	return func(c *gin.Context) {
	}
}

// Creates handler func for deleting sheet
func NewSheetDeleteHandler(sheetManager app.ISheetManager) func(c *gin.Context) {
	return func(c *gin.Context) {
	}
}

// Creates handler func for updating sheet
func NewSheetUpdateHanlder(sheetManager app.ISheetManager) func(c *gin.Context) {
	return func(c *gin.Context) {
	}
}
