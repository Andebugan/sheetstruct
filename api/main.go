package main

import (
	"context"
	"log"

	"github.com/andebugan/sheetstruct/docs"
	"github.com/andebugan/sheetstruct/internal/api"
	"github.com/andebugan/sheetstruct/internal/app/managers"
	"github.com/andebugan/sheetstruct/internal/db/dbs"
	"github.com/andebugan/sheetstruct/internal/db/repos"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// To regenerate swagger documentation:
// 1. install swag via `go install github.com/swaggo/swag/cmd/swag@latest`
// 2. run `swag init`
// More info in swag repo: https://github.com/swaggo/swag
// To access api go to http://host:port/swagger/index.html
//
//	@securityDefinitions.apikey  BearerAuth
//	@in                          header
//	@name                        Authorization
//	@description                 Type "Bearer" followed by a space and JWT token.
func main() {
	docs.SwaggerInfo.Title = "Sheetstruct Player List API"
	docs.SwaggerInfo.Description = "API for Sheetstruct Player List Management Server"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api/v1"

	// initialize database
	db, err := dbs.InitSqliteDatabase()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return
	}
	defer dbs.CloseSqliteDatabase(db)

	// Initalize repositories
	userRepository := repos.NewGormUserRepository(db)
	sheetRepository := repos.NewGormSheetRepository(db)
	componentRepository := repos.NewGormComponentRepository(db)
	ctx := context.Background()

	// Initialize managers
	userManager := managers.NewUserManager(ctx, userRepository)
	sheetManager := managers.NewSheetManager(ctx, userRepository, sheetRepository)
	componentManager := managers.NewComponentManager(ctx, userRepository, sheetRepository, componentRepository)

	// Initialize api
	router := gin.Default()
	router.Use(api.CorsMiddleware())

	v1 := router.Group("/api/v1")
	{
		// health check handler
		v1.GET("/health", api.HealthCheck)

		// user handlers
		v1.POST("/user", api.NewUserCreateHandler(userManager))
		v1.POST("/user/login", api.NewUserLoginHandler(userManager))
		v1.POST("/user/refresh_token", api.NewUserRefreshTokenHandler(userManager))
	}

	v1_protected := router.Group("/api/v1").Use(api.JwtMiddleware())
	{
		// user handlers
		v1_protected.GET("/user", api.NewUserGetCurrentHandler(userManager))
		v1_protected.PATCH("/user", api.NewUserUpdateHandler(userManager))
		v1_protected.DELETE("/user", api.NewUserDeleteHandler(userManager))

		// sheet handlers
		v1_protected.GET("/sheets", api.NewSheetsGetHandler(sheetManager))
		v1_protected.GET("/sheet/:sid", api.NewSheetGetHandler(sheetManager))
		v1_protected.POST("/sheet", api.NewSheetCreateHandler(sheetManager))
		v1_protected.PATCH("/sheet", api.NewSheetUpdateHanlder(sheetManager))
		v1_protected.PUT("/sheet/:sid", api.NewSheetCloneHanlder(sheetManager))
		v1_protected.DELETE("/sheet/:sid", api.NewSheetDeleteHandler(sheetManager))

		// component handlers
		v1_protected.GET("/sheet/:sid/component/:cid", api.NewComponentGetHandler(componentManager))
		v1_protected.GET("/sheet/:sid/components", api.NewComponentsGetHandler(componentManager))
		v1_protected.PUT("/sheet/:sid/component/:cid", api.NewComponentCloneHandler(componentManager))
		v1_protected.POST("/sheet/:sid/component", api.NewComponentCreateHandler(componentManager))
		v1_protected.PATCH("/sheet/:sid/component", api.NewComponentUpdateHandler(componentManager))
		v1_protected.DELETE("/sheet/:sid/component/:cid", api.NewComponentDeleteHandler(componentManager))
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run("0.0.0.0:5000")
}
