package api

import (
	"github.com/andebugan/sheetstruct/internal/app"
	"github.com/gin-gonic/gin"
)

// Creates handler func for getting component
func NewComponentGetHandler(componentManager app.IComponentManager) func(c *gin.Context) {
	return func(c *gin.Context) {
	}
}

// Creates handler func for getting many components
func NewComponentsGetHandler(componentManager app.IComponentManager) func(c *gin.Context) {
	return func(c *gin.Context) {
	}
}

// Creates handler func for cloning components
func NewComponentCloneHandler(componentManager app.IComponentManager) func(c *gin.Context) {
	return func(c *gin.Context) {
	}
}

// Creates handler func for creating new component
func NewComponentCreateHandler(componentManager app.IComponentManager) func(c *gin.Context) {
	return func(c *gin.Context) {
	}
}

// Creates handler func for deleting component
func NewComponentDeleteHandler(componentManager app.IComponentManager) func(c *gin.Context) {
	return func(c *gin.Context) {
	}
}

// Creates handler func for updating compomnent
func NewComponentUpdateHandler(componentManager app.IComponentManager) func(c *gin.Context) {
	return func(c *gin.Context) {
	}
}
