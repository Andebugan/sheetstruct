package api

import (
	"fmt"
	"net/http"

	"github.com/andebugan/sheetstruct/internal/app"
	"github.com/andebugan/sheetstruct/internal/app/managers"
	"github.com/andebugan/sheetstruct/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

// @description DTO for recieving generated JWT token
type TokenDTO struct {
	Token        string
	RefreshToken string
}

// Creates endpoint for new user generation
//
//	@summary		Creates new user
//	@description 	Creates new user from name, email and password, when provided with correct description ans password
//	@tags			User
//	@accept 		json
//	@produce 		json
//	@success 		201 {object} models.User
//	@param			request body UserCredsDTO true "Credentials DTO"
//	@failure		400 {object} string
//	@failure		409 {object} string
//	@failure		500 {object} string
//	@router 		/user [post]
func NewUserCreateHandler(userManager managers.IUserManager) func(c *gin.Context) {
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

		if err == app.ErrUserNameAlredyExists || err == app.ErrUserEmailAlredyExists {
			c.JSON(http.StatusConflict, err.Error())
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusCreated, user)
	}
}

// Get current user endpoint
//
//	@summary		Gets current authenticated User
//	@description 	Get user information for authenticated user via BearerAuth token
//	@tags			User
//	@produce 		json
//	@success 		200 {object} models.User
//	@failure		401 {object} string
//	@failure		404 {object} string
//	@failure		500 {object} string
//	@security       BearerAuth
//	@router 		/user [get]
func NewUserGetCurrentHandler(userManager managers.IUserManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		uid, err := TryGetUidFromToken(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		fmt.Println(uid)
		user, err := userManager.Get(uid)

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
//
//	@summary 		Delete current user
//	@description 	Deletes currently authenticated User from id recieved via BearerAuth token
//	@tags			User
//	@success 		200 {object} string
//	@failure		401 {object} string
//	@failure		404 {object} string
//	@failure		500 {object} string
//	@security       BearerAuth
//	@router 	 	/user [delete]
func NewUserDeleteHandler(userManager managers.IUserManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		uid, err := TryGetUidFromToken(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		err = userManager.Delete(uid)

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
//
//	@summary 		Update user
//	@description 	Updates User with values from recieved user object
//	@tags			User
//	@param			request body models.User true "New user data"
//	@success 		200 {object} models.User
//	@failure 		401 {object} string
//	@failure 		404 {object} string
//	@failure 		500 {object} string
//	@security       BearerAuth
//	@router 		/user [patch]
func NewUserUpdateHandler(userManager managers.IUserManager) func(c *gin.Context) {
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
//
//	@summary 		User login
//	@description 	Authenticates user and generates access token
//	@tags			User
//	@param 			request body UserAuthDTO true "User credentials"
//	@accept			json
//	@produce		json
//	@success 		200 {object} TokenDTO
//	@failure		400 {object} string
//	@failure		400 {object} string
//	@failure		500 {object} string
//	@router 		/user/login [post]
func NewUserLoginHandler(userManager managers.IUserManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		var authData UserAuthDTO
		err := c.BindJSON(&authData)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		user, err := userManager.Login(authData.Login, authData.Password)

		if err == app.ErrUserAuthFailed {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		token, err := GenerateJWT(user.Id, user.Name)
		refreshToken, err := GenerateRefreshToken(user.Id, user.Name)

		c.JSON(http.StatusOK, TokenDTO{
			Token:        token,
			RefreshToken: refreshToken,
		})
	}
}

// Creates handler func for refreshing auth token
//
//	@summary 		Refresh auth token
//	@description 	Refreshes expired JWT token
//	@tags			User
//	@success 		200 {string} token
//	@failure		400 {object} string
//	@failure		401 {object} string
//	@failure		500 {object} string
//	@router 		/user/refresh_token [post]
func NewUserRefreshTokenHandler(userManager managers.IUserManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req struct {
			RefreshToken string `json:"refresh_token"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		token, err := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (any, error) {
			return RefreshSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, "Invalid refresh token")
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		var idValue uint64
		switch v := claims["id"].(type) {
		case uint64:
			idValue = v
		case float64:
			idValue = uint64(v)
		case int:
			idValue = uint64(v)
		case int64:
			idValue = uint64(v)
		default:
			c.JSON(http.StatusBadRequest, "Invalid token claims")
			return
		}
		id := models.UserID{
			Value: idValue,
		}

		user, err := userManager.Get(id)
		if err == app.ErrUserNotFound {
			c.JSON(http.StatusNotFound, err.Error())
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		newToken, err := GenerateJWT(user.Id, user.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, newToken)
	}
}
