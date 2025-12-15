package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// @description	Health check response structure
type HealthResponse struct {
	Status string `json:"status" example:"healthy"`
}

// Healthcheck handle
//
//	@summary		Health check endpoint
//	@description	Checks if the server is active
//  @tags			Utility
//	@accept			json
//	@produce		json
//	@success		200	{object} HealthResponse
//	@router			/health [get]
func healthCheck(c *gin.Context) {
	responce := HealthResponse{
		Status: "healty",
	}

	c.JSON(http.StatusOK, responce)
}
