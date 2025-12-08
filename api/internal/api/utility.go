package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @description	Health check response structure
type HealthResponse struct {
	Status string `json:"status" example:"healthy"`
}

// Healthcheck handle
//
//	@summary		Health check endpoint
//	@description	Checks if the server is active
//	@accept			json
//	@produce		json
//	@success		200	{object}	HealthResponse
//	@router			/health [get]
func healthCheck(c *gin.Context) {
	responce := HealthResponse{
		Status: "healty",
	}

	c.JSON(http.StatusOK, responce)
}
