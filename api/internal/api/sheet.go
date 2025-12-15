package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/andebugan/sheetstruct/internal/app"
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
//	@params			request body models.SheetFilter true "Filter settings"
//  @success 		200 {object} []app.SheetFilter
//  @failure		401
//  @failure		500
//  @router 		/sheets [get]
func NewSheetsGetHandler(sheetManager app.ISheetManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		uid, err := TryGetUidFromToken(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		// Get filter
		var filter app.SheetFilter
		err = c.BindJSON(&filter)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		
		sheets, err := sheetManager.GetMany(filter, uid)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, sheets)
	}
}

// Creates handler func for getting single sheet
//
//  @summary		Get sheets
//  @description 	Get multiple sheets avaliable to user
//  @tags			Sheet
//  @accept 		json
//  @produce 		json
//	@param 			sid query int true "Sheet ID"
//  @success 		200 {object} models.Sheet
//  @failure		401
//  @failure		404
//  @failure		500
//  @router 		/sheet/{sid} [get]
func NewSheetGetHandler(sheetManager app.ISheetManager) func(c *gin.Context) {
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
//  @failure 		401
//  @failure		500
//  @router 		/sheet [post]
func NewSheetCreateHandler(sheetManager app.ISheetManager) func(c *gin.Context) {
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
//  @success 		200
//  @failure		401
//  @failure		404
//  @failure		500
//  @router 		/sheet/{sid} [delete]
func NewSheetDeleteHandler(sheetManager app.ISheetManager) func(c *gin.Context) {
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
//  @description 	Updates sheet by id
//  @tags			Sheet
//  @accept 		json
//  @produce 		json
//	@params			request body models.Sheet true "Updated sheet"
//  @success 		200 {object} models.Sheet
//  @failure		401
//  @failure		404
//  @failure		500
//  @router 		/sheet/{sid} [patch]
func NewSheetUpdateHanlder(sheetManager app.ISheetManager) func(c *gin.Context) {
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
//  @failure		401
//  @failure		404
//  @failure		500
//  @router 		/sheet/{sid} [put]
func NewSheetCloneHanlder(sheetManager app.ISheetManager) func(c *gin.Context) {
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
