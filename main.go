package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/natha/kyoyu-backend/controllers"
	"github.com/natha/kyoyu-backend/initializers"
	"github.com/natha/kyoyu-backend/middleware"
	"github.com/natha/kyoyu-backend/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()

	// Configure CORS middleware to allow requests from all origins
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowHeaders = []string{"Authorization", "Content-Type", "x-xsrf-token"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	// Routes
	api := r.Group("/api")
	{
		api.POST("/signup", controllers.SignUp)
		api.POST("/login", controllers.Login)
		api.POST("/logout", controllers.Logout)

		api.GET("/user", func(c *gin.Context) {
			middleware.RequireAuth(c, "")
		}, controllers.UserShow)

		my := api.Group("/my")
		{
			my.Use(func(c *gin.Context) {
				middleware.RequireAuth(c, models.UserRole("member"))
			})
			my.GET("/posts", controllers.PostIndex)
			my.GET("/posts/mine", controllers.PostShow)
			my.POST("/posts", controllers.PostCreate)
		}

		admin := api.Group("/admin")
		{
			admin.Use(func(c *gin.Context) {
				middleware.RequireAuth(c, models.UserRole("admin"))
			})
			admin.GET("/users", controllers.UserIndex)
		}
	}

	r.Run()
}
