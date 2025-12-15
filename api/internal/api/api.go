package api

import (
	"github.com/andebugan/sheetstruct/internal/app"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: get actual source from env config
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Set's up and returns router for HTTP REST API
// To regenerate swagger documentation:
// 1. install swag via `go install github.com/swaggo/swag/cmd/swag@latest`
// 2. run `swag init`
// More info in swag repo: https://github.com/swaggo/swag
// To access api go to http://host:port/swagger/index.html
func Setup() *gin.Engine {
	userManager := app.NewUserManager()
	sheetManager := app.NewSheetManager()
	componentManager := app.NewComponentManager()

	router := gin.Default()
	router.Use(corsMiddleware())

	v1 := router.Group("/api/v1")
	{
		// health check handler
		v1.GET("/health", healthCheck)

		// user handlers
		v1.POST("/user", NewUserCreateHandler(userManager))
		v1.POST("/user/login", NewUserLoginHandler(userManager))
		v1.POST("/user/refresh_token", NewUserRefreshTokenHandler(userManager))
	}

	v1_protected := v1.Group("/api/v1").Use(JwtMiddleware())
	{
		// user handlers
		v1_protected.GET("/user", NewUserGetCurrentHandler(userManager))
		v1_protected.PATCH("/user", NewUserUpdateHandler(userManager))
		v1_protected.DELETE("/user", NewUserDeleteHandler(userManager))
		v1_protected.POST("/user/logout", NewUserLogoutHandler(userManager))

		// sheet handlers
		v1_protected.GET("/sheets", NewSheetsGetHandler(sheetManager))
		v1_protected.GET("/sheet/:sid", NewSheetGetHandler(sheetManager))
		v1_protected.POST("/sheet", NewSheetCreateHandler(sheetManager))
		v1_protected.PATCH("/sheet", NewSheetUpdateHanlder(sheetManager))
		v1_protected.PUT("/sheet/:sid", NewSheetCloneHanlder(sheetManager))
		v1_protected.DELETE("/sheet/:sid", NewSheetDeleteHandler(sheetManager))

		// component handlers
		v1_protected.GET("/sheet/:sid/component/:cid", NewComponentGetHandler(componentManager))
		v1_protected.GET("/sheet/:sid/components", NewComponentsGetHandler(componentManager))
		v1_protected.PUT("/sheet/:sid/component/:cid", NewComponentCloneHandler(componentManager))
		v1_protected.POST("/sheet/:sid/component", NewComponentCreateHandler(componentManager))
		v1_protected.PATCH("/sheet/:sid/component", NewComponentUpdateHandler(componentManager))
		v1_protected.DELETE("/sheet/:sid/component/:cid", NewComponentDeleteHandler(componentManager))
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
