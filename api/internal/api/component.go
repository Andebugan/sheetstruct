package api

import (
	"net/http"

	"github.com/andebugan/sheetstruct/internal/app"
	"github.com/andebugan/sheetstruct/internal/models"
	"github.com/gin-gonic/gin"
)

// Creates handler func for getting component
func NewComponentGetHandler(componentManager app.IComponentManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		// TODO: replace with ID from JWT Token
		uid := models.UserID{}
		// TODO: replace with ID from url
		sid := models.SheetID{}

		components, err := componentManager.GetMany(sid, uid)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, components)
	}
}

// Creates handler func for getting many components
func NewComponentsGetHandler(componentManager app.IComponentManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		// TODO: replace with ID from JWT Token
		uid := models.UserID{}
		// TODO: replace with ID from request params
		sid := models.SheetID{}
		// TODO: replace with ID from request params
		cid := models.ComponentID{}

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
func NewComponentCloneHandler(componentManager app.IComponentManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		// TODO: replace with ID from JWT Token
		uid := models.UserID{}
		// TODO: replace with ID from request params
		sid := models.SheetID{}
		// TODO: replace with ID from request params
		cid := models.ComponentID{}

		component, err := componentManager.Clone(cid, sid, uid)

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

// Creates handler func for creating new component
func NewComponentCreateHandler(componentManager app.IComponentManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		// TODO: replace with ID from JWT Token
		uid := models.UserID{}
		// TODO: replace with ID from request params
		sid := models.SheetID{}

		component, err := componentManager.Create(sid, uid)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, component)
	}
}

// Creates handler func for deleting component
func NewComponentDeleteHandler(componentManager app.IComponentManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		// TODO: replace with ID from JWT Token
		uid := models.UserID{}
		// TODO: replace with ID from request params
		sid := models.SheetID{}
		// TODO: replace with ID from request params
		cid := models.ComponentID{}

		err := componentManager.Delete(cid, sid, uid)

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
func NewComponentUpdateHandler(componentManager app.IComponentManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		var component models.Component
		err := c.BindJSON(&component)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		updatedComponent, err := componentManager.Update(component)

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
