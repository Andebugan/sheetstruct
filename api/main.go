package main

import (
	"github.com/andebugan/sheetstruct/docs"
	"github.com/andebugan/sheetstruct/internal/api"
)

// To regenerate swagger documentation:
// 1. install swag via `go install github.com/swaggo/swag/cmd/swag@latest`
// 2. run `swag init`
// More info in swag repo: https://github.com/swaggo/swag
// To access api go to http://host:port/swagger/index.html
//  @securityDefinitions.apikey  BearerAuth
//  @in                          header
//  @name                        Authorization
//  @description                 Type "Bearer" followed by a space and JWT token.
func main() {
	docs.SwaggerInfo.Title = "Sheetstruct Player List API"
	docs.SwaggerInfo.Description = "API for Sheetstruct Player List Management Server"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api/v1"

	router := api.Setup()
	router.Run("0.0.0.0:5000")
}
