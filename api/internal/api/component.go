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

var ErrComponentIdParsingFailed = errors.New("Unable to parse ComponentID from query")

func TryGetComponentIdFromQuery(c *gin.Context) (models.ComponentID, error) {
	cid := models.ComponentID{}

	cidRaw, err := strconv.Atoi(c.Query("cid"))
	if err != nil {
		return cid, ErrComponentIdParsingFailed
	}

	cid.Value = uint64(cidRaw)

	return cid, nil
}

// Creates handler func for getting component
//
//  @summary		Gets components
//  @description 	Creates new component inside specified sheet
//  @tags			Component
//  @accept 		json
//  @produce 		json
//	@param			sid query int true "Sheet ID"
//  @success 		200 {object} models.Component
//  @failure		401 {object} string
//  @failure		500 {object} string
//  @router 		/sheet/{sid}/component [get]
//	@security       BearerAuth
func NewComponentsGetHandler(componentManager managers.IComponentManager) func(c *gin.Context) {
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

		components, err := componentManager.GetMany(sid, uid)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, components)
	}
}

// Creates handler func for getting single component
//
//  @summary		Get component
//  @description 	Get gets single component inside specified sheet
//  @tags			Component
//  @accept 		json
//  @produce 		json
//	@param			sid query int true "Sheet ID"
//	@param			cid query int true "Component ID"
//  @success 		200 {object} []models.Component
//  @failure		404 {object} string
//  @failure		500 {object} string
//  @router 		/sheet/{sid}/component/{cid} [get]
//	@security       BearerAuth
func NewComponentGetHandler(componentManager managers.IComponentManager) func(c *gin.Context) {
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

		cid, err := TryGetComponentIdFromQuery(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		component, err := componentManager.Get(cid, sid, uid)

		if err == app.ErrComponentNotFound {
			c.JSON(http.StatusNotFound, err.Error())
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, component)
	}
}

// Creates handler func for cloning components
//
//  @summary		Clone component
//  @description 	Creates copy of component inside specified sheet
//  @tags			Component
//  @accept 		json
//  @produce 		json
//	@param			sid query int true "Sheet ID"
//	@param			cid query int true "Component ID"
//  @success 		201 {object} models.Component
//  @failure		401 {object} string
//  @failure		404 {object} string
//  @failure		500 {object} string
//  @router 		/sheet/{sid}/component/{cid} [put]
//	@security       BearerAuth
func NewComponentCloneHandler(componentManager managers.IComponentManager) func(c *gin.Context) {
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

		cid, err := TryGetComponentIdFromQuery(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		component, err := componentManager.Clone(cid, sid, uid)

		if err == app.ErrComponentNotFound {
			c.JSON(http.StatusNotFound, err.Error())
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusCreated, component)
	}
}

// Creates handler func for creating new component
//
//  @summary		Create new component
//  @description 	Creates new component inside specified sheet
//  @tags			Component
//  @accept 		json
//  @produce 		json
//	@param			sid query int true "Sheet ID"
//  @success 		201 {object} models.Component
//  @failure		401 {object} string
//  @failure		500 {object} string
//  @router 		/sheet/{sid}/component [post]
//	@security       BearerAuth
func NewComponentCreateHandler(componentManager managers.IComponentManager) func(c *gin.Context) {
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

		component, err := componentManager.Create(sid, uid)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusCreated, component)
	}
}

// Creates handler func for deleting component
//
//  @summary		Delete component
//  @description 	Deletes component inside specified sheet
//  @tags			Component
//  @accept 		json
//  @produce 		json
//  @success 		200 {object} string
//	@param			sid query int true "Sheet ID"
//	@param			cid query int true "Component ID"
//  @failure		404 {object} string
//  @failure		500 {object} string
//  @router 		/sheet/{sid}/component/{cid} [delete]
//	@security       BearerAuth
func NewComponentDeleteHandler(componentManager managers.IComponentManager) func(c *gin.Context) {
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

		cid, err := TryGetComponentIdFromQuery(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		err = componentManager.Delete(cid, sid, uid)

		if err == app.ErrComponentNotFound {
			c.JSON(http.StatusNotFound, err.Error())
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, "Component successfully deleted")
	}
}

// Creates handler func for updating component
//
//  @summary		Updates component
//  @description 	Updates component with recieved data
//  @tags			Component
//  @accept 		json
//  @produce 		json
//  @success 		200 {object} models.Component
//	@param			request body models.Component true "Updated component"
//  @failure		400 {object} string
//  @failure		401 {object} string
//  @failure		404 {object} string
//  @failure		500 {object} string
//  @router 		/sheet/{sid}/component [patch]
//	@security       BearerAuth
func NewComponentUpdateHandler(componentManager managers.IComponentManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		uid, err := TryGetUidFromToken(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		var component models.Component
		err = c.BindJSON(&component)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		updatedComponent, err := componentManager.Update(component, uid)

		if err == app.ErrComponentNotFound {
			c.JSON(http.StatusNotFound, err.Error())
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, updatedComponent)
	}
}
