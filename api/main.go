package main

import (
	"github.com/andebugan/sheetstruct/internal/api"
	"github.com/andebugan/sheetstruct/docs"
)

// To regenerate swagger documentation:
// 1. install swag via `go install github.com/swaggo/swag/cmd/swag@latest`
// 2. run `swag init`
// More info in swag repo: https://github.com/swaggo/swag

func main() {
	docs.SwaggerInfo.Title = "Sheetstruct Player List API"
	docs.SwaggerInfo.Description = "API for Sheetstruct Player List Management Server"
	docs.SwaggerInfo.Version = "1.0"

	router := api.Setup()
	router.Run(":3010")
}
