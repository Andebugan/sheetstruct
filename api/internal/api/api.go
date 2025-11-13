package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

// Set's up and returns router for HTTP REST API
func Setup() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", healthCheck)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}

// Health check responce model
// @description Health check response structure
type HealthResponse struct {
	Status string `json:"status" example:"healthy"`
}

// Healthcheck handle
// @summary Health check endpoint
// @description Checks if the server is active
// @accept json
// @produce json
// @success 200 {object} HealthResponse
// @router /health [get]
func healthCheck(c *gin.Context) {
	responce := HealthResponse {
		Status: "healty",
	}

	c.JSON(http.StatusOK, responce)
}


