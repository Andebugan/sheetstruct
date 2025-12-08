package api

import (
	"net/http"

	"github.com/andebugan/sheetstruct/internal/app"
	"github.com/andebugan/sheetstruct/internal/models"
	"github.com/gin-gonic/gin"
)

// @description DTO for new user info
type UserCredsDTO struct {
	Name     string
	Email    string
	Password string
}

// @description DTO for authenticating new user
type UserAuthDTO struct {
	Login    string
	Password string
}

// Creates new user from name, email and password
// @summary User creation endpoint
// @description Creates new user, when provided with correct description ans password
// @accept json
// @produce json
// @success 200 {object} models.User
// @router /user [post]
func NewUserCreateHandler(userManager app.IUserManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		var newUser UserCredsDTO
		err := c.BindJSON(&newUser)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		var userData = models.NewUserData{
			Name:     newUser.Name,
			Email:    newUser.Email,
			Password: newUser.Password,
		}
		user, err := userManager.Create(userData)

		if err == app.ErrUserAlredyExists {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

// Creates handler func for getting current User info
func NewUserGetCurrentHandler(userManager app.IUserManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		// TODO: replace with ID from JWT Token
		id := models.UserID{}

		user, err := userManager.Get(id)

		if err == app.ErrUserNotFound {
			c.JSON(http.StatusNotFound, err.Error())
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

// Creates handler func for deleting User
func NewUserDeleteHandler(userManager app.IUserManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		// TODO: replace with ID from JWT Token
		id := models.UserID{}

		err := userManager.Delete(id)

		if err == app.ErrUserNotFound {
			c.JSON(http.StatusNotFound, err.Error())
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, "User deleted successfully")
	}
}

// Creates handler func for updating User
func NewUserUpdateHandler(userManager app.IUserManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		var user models.User
		err := c.BindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		updatedUser, err := userManager.Update(user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, updatedUser)
	}
}

// Creates handler func for authenticating User
func NewUserAuthHandler(userManager app.IUserManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		var authData UserAuthDTO
		err := c.BindJSON(&authData)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		user, err := userManager.Auth(authData.Login, authData.Password)

		if err == app.ErrUserAuthFailed {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
