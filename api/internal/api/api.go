package api

import (
	"github.com/andebugan/sheetstruct/internal/app"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

// Set's up and returns router for HTTP REST API
func Setup() *gin.Engine {
	userManager := app.NewUserManager()
	sheetManager := app.NewSheetManager()
	componentManager := app.NewComponentManager()

	router := gin.Default()
	v1 := router.Group("/api/v1")
	{
		// health check handler
		v1.GET("/health", healthCheck)

		// user handlers
		v1.GET("/user", NewUserGetHandler(userManager))
		v1.POST("/user", NewUserCreateHandler(userManager))
		v1.PATCH("/user", NewUserUpdateHandler(userManager))
		v1.DELETE("/user", NewUserDeleteHandler(userManager))
		v1.POST("/user/auth", NewUserAuthHandler(userManager))
		v1.POST("/user/reset", NewUserResetHandler(userManager))

		// sheet handlers
		v1.GET("/sheet", NewSheetGetHandler(sheetManager))
		v1.GET("/sheets", NewSheetsGetHandler(sheetManager))
		v1.POST("/sheet", NewSheetCreateHandler(sheetManager))
		v1.PATCH("/sheet", NewSheetUpdateHanlder(sheetManager))
		v1.DELETE("/sheet", NewSheetDeleteHandler(sheetManager))

		// component handlers
		v1.GET("/component", NewComponentGetHandler(componentManager))
		v1.PUT("/component/clone", NewComponentCloneHandler(componentManager))
		v1.POST("/component", NewComponentCreateHandler(componentManager))
		v1.PATCH("/component", NewComponentUpdateHandler(componentManager))
		v1.DELETE("/component", NewComponentDeleteHandler(componentManager))
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
