package api

import (
	"net/http"

	"github.com/andebugan/sheetstruct/internal/app"
	"github.com/andebugan/sheetstruct/internal/models"
	"github.com/gin-gonic/gin"
)

// Creates handler func for getting multiple sheets
func NewSheetsGetHandler(sheetManager app.ISheetManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		// TODO: replace with ID from JWT Token
		uid := models.UserID{}
		
		sheets, err := sheetManager.GetMany(uid)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, sheets)
	}
}

// Creates handler func for getting single sheet
func NewSheetGetHandler(sheetManager app.ISheetManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		// TODO: replace with ID from JWT Token
		uid := models.UserID{}
		// TODO: replace with ID from request params
		sid := models.SheetID{}
		
		sheet, err := sheetManager.Get(sid, uid)

		if err == app.ErrUserNotFound {
			c.JSON(http.StatusNotFound, err.Error())
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, sheet)
	}
}

// Creates handler func for creating new sheet
func NewSheetCreateHandler(sheetManager app.ISheetManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		// TODO: replace with ID from JWT Token
		uid := models.UserID{}
		
		sheet, err := sheetManager.Create(uid)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, sheet)
	}
}

// Creates handler func for deleting sheet
func NewSheetDeleteHandler(sheetManager app.ISheetManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		// TODO: replace with ID from JWT Token
		uid := models.UserID{}
		// TODO: replace with ID from request params
		sid := models.SheetID{}
		
		err := sheetManager.Delete(sid, uid)

		if err == app.ErrUserNotFound {
			c.JSON(http.StatusNotFound, err.Error())
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, "Sheet successfully deleted")
	}
}

// Creates handler func for updating sheet
func NewSheetUpdateHanlder(sheetManager app.ISheetManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		var sheet models.Sheet
		err := c.BindJSON(&sheet)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		updatedSheet, err := sheetManager.Update(sheet)

		if err == app.ErrUserNotFound {
			c.JSON(http.StatusNotFound, err.Error())
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, updatedSheet)
	}
}
