package router

import (
	"github.com/gin-gonic/gin"
	"os"
	"task5-btpns-stella/controllers"
	docs "task5-btpns-stella/docs"
	"task5-btpns-stella/middlewares"
)

func WebRoute(router *gin.Engine) {
	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := router.Group("/api/v1")
	{
		user := v1.Group("/user")
		{
			user.GET("/login", controllers.LoginController)
			user.POST("/register", controllers.RegisterController)
			user.PUT("/update-user", middlewares.TokenAuthMiddleware(), controllers.UserUpdateController)
			user.DELETE("/delete-user", middlewares.TokenAuthMiddleware(), controllers.UserDeleteController)
			user.POST("/check-session", middlewares.TokenAuthMiddleware(), controllers.CheckSessionController)
		}
		photo := v1.Group("/photo")
		{
			photo.POST("/photos", middlewares.TokenAuthMiddleware(), controllers.PhotoPostController)
			photo.GET("/photos", middlewares.TokenAuthMiddleware(), controllers.PhotoGetController)
			photo.PUT("/photoid", middlewares.TokenAuthMiddleware(), controllers.PhotoEditController)
			photo.DELETE("/photoid", middlewares.TokenAuthMiddleware(), controllers.PhotoDeleteController)
		}
	}
}

func Route() (*gin.Engine, string) {
	appRelease := os.Getenv("APP_MODE")
	if appRelease == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(CORS())

	//middlewares.InitLogServer(nil)
	WebRoute(router)
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "9191"
	}

	//attachSwagger(router, port)
	// log.Fatal(router.Run(":" + port))
	return router, port
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		//c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
