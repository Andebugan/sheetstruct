package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/andebugan/sheetstruct/internal/app"
	"github.com/andebugan/sheetstruct/internal/app/managers"
	"github.com/andebugan/sheetstruct/internal/models"
	"github.com/gin-gonic/gin"
)

var ErrSheetIdParsingFailed = errors.New("Unable to parse SheetID from query")

func TryGetSheetIdFromQuery(c *gin.Context) (models.SheetID, error) {
	sid := models.SheetID{}

	sidRaw, err := strconv.Atoi(c.Query("sid"))
	if err != nil {
		return sid, ErrSheetIdParsingFailed
	}

	sid.Value = uint64(sidRaw)

	return sid, nil
}

// Creates handler func for getting multiple sheets
//
//  @summary		Get sheets
//  @description 	Returns filtered sheet collection avaliable to user
//  @tags			Sheet
//  @accept 		json
//  @produce 		json
//  @success 		200 {object} []models.Sheet
//  @failure		401 {object} string
//  @failure		500 {object} string
//  @router 		/sheets [get]
//	@security       BearerAuth
func NewSheetsGetHandler(sheetManager managers.ISheetManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		uid, err := TryGetUidFromToken(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		sheets, err := sheetManager.GetMany(uid)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, sheets)
	}
}

// Creates handler func for getting single sheet
//
//  @summary		Get sheet by id
//  @description 	Get single sheet by id
//  @tags			Sheet
//  @accept 		json
//  @produce 		json
//	@param 			sid query int true "Sheet ID"
//  @success 		200 {object} models.Sheet
//  @failure		401 {object} string
//  @failure		404 {object} string
//  @failure		500 {object} string
//  @router 		/sheet/{sid} [get]
//	@security       BearerAuth
func NewSheetGetHandler(sheetManager managers.ISheetManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		uid, err := TryGetUidFromToken(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		sid, err := TryGetSheetIdFromQuery(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}		

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
//
//  @summary		Create new sheet
//  @description 	Creates new empty sheet
//  @tags			Sheet
//  @accept 		json
//  @produce 		json
//  @success 		201 {object} models.Sheet
//  @failure 		401 {object} string
//  @failure		500 {object} string
//  @router 		/sheet [post]
//	@security       BearerAuth
func NewSheetCreateHandler(sheetManager managers.ISheetManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		uid, err := TryGetUidFromToken(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
				
		sheet, err := sheetManager.Create(uid)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusCreated, sheet)
	}
}

// Creates handler func for deleting sheet
//
//  @summary		Delete sheet
//  @description 	Deletes sheet by id
//  @tags			Sheet
//  @accept 		json
//  @produce 		json
//	@param 			sid query int true "Sheet ID"
//  @success 		200 {object} string
//  @failure		401 {object} string
//  @failure		404 {object} string
//  @failure		500 {object} string
//  @router 		/sheet/{sid} [delete]
//	@security       BearerAuth
func NewSheetDeleteHandler(sheetManager managers.ISheetManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		uid, err := TryGetUidFromToken(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		sid, err := TryGetSheetIdFromQuery(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}		

		err = sheetManager.Delete(sid, uid)

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
//
//  @summary		Updates sheet
//  @description 	Updates sheet
//  @tags			Sheet
//  @accept 		json
//  @produce 		json
//	@param			request body models.Sheet true "Updated sheet"
//  @success 		200 {object} models.Sheet
//  @failure		401 {object} string
//  @failure		404 {object} string
//  @failure		500 {object} string
//  @router 		/sheet [patch]
//	@security       BearerAuth
func NewSheetUpdateHanlder(sheetManager managers.ISheetManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		uid, err := TryGetUidFromToken(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		var sheet models.Sheet
		err = c.BindJSON(&sheet)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		updatedSheet, err := sheetManager.Update(sheet, uid)

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

// Creates handler func for cloning sheet
//
//  @summary		Clone sheet
//  @description 	Create deep clone for sheet with provided sid
//  @tags			Sheet
//  @accept 		json
//  @produce 		json
//	@param			sid query int true "Sheet ID"
//  @success 		201 {object} models.Sheet
//  @failure		401 {object} string
//  @failure		404 {object} string
//  @failure		500 {object} string
//  @router 		/sheet/{sid} [put]
//	@security       BearerAuth
func NewSheetCloneHanlder(sheetManager managers.ISheetManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		uid, err := TryGetUidFromToken(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		sidRaw, err := strconv.Atoi(c.Param("sid"))
		if err != nil {
			c.JSON(http.StatusBadRequest, "Invalid sheet ID")
			return
		}
		
		sid := models.SheetID{Value: uint64(sidRaw)}

		sheet, err := sheetManager.Clone(sid, uid)

		if err == app.ErrUserNotFound {
			c.JSON(http.StatusNotFound, err.Error())
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusCreated, sheet)
	}
}
