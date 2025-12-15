package api

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/andebugan/sheetstruct/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// TODO: get value from env config
var jwtKey = []byte("jwt_key")
// TODO: get value from env config
var RefreshSecret = []byte("jwt_refresh_secret")

var ErrUnableToGetUidFromToken = errors.New("Unable to get uid from JWT token")

func TryGetUidFromToken(c *gin.Context) (models.UserID, error) {
	value, exists := c.Get("id")

	uid := models.UserID{}

	if !exists {
		return uid, ErrUnableToGetUidFromToken
	} else {
		uid.Value = value.(uint64)
	}

	return uid, nil
}

func GenerateJWT(id models.UserID, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       id.Value,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	})
	return token.SignedString(jwtKey)
}

func GenerateRefreshToken(id models.UserID, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       id.Value,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
	})
	return token.SignedString(RefreshSecret)
}

func JwtMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, "Missing token")
            c.Abort()
            return
        }

        tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

        token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
            return jwtKey, nil
        })

        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            c.Set("id", claims["id"])
            c.Set("username", claims["username"])
            c.Next()
        } else {
            c.JSON(http.StatusUnauthorized, "Invalid token")
            c.Abort()
        }
    }
} 
