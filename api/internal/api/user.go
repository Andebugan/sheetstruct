package api

import (
	"github.com/andebugan/sheetstruct/internal/app"
	"github.com/gin-gonic/gin"
)

// Creates handler func for creating User
func NewUserCreateHandler(userManager app.IUserManager) func(c *gin.Context) {
	// Creates new user from name, email and password
	// @summary User creation endpoint
	// @description Creates new user, when provided with correct description ans password
	// @accept json
	// @produce json
	// @success 200 {object} User
	// @router /user [post]
	return func(c *gin.Context) {
	}
}

// Creates handler func for getting User info 
func NewUserGetHandler(userManager app.IUserManager) func(c *gin.Context) {
	return func(c *gin.Context) {
	}
}

// Creates handler func for deleting User
func NewUserDeleteHandler(userManager app.IUserManager) func(c *gin.Context) {
	return func(c *gin.Context) {
	}
}

// Creates handler func for updating User
func NewUserUpdateHandler(userManager app.IUserManager) func(c *gin.Context) {
	return func(c *gin.Context) {
	}
}

// Creates handler func for authenticating User
func NewUserAuthHandler(userManager app.IUserManager) func(c *gin.Context) {
	return func(c *gin.Context) {
	}
}

// Creates handler func for resetting User password
func NewUserResetHandler(userManager app.IUserManager) func(c *gin.Context) {
	return func(c *gin.Context) {
	}
}
