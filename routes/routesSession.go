package routes

import (
	"session-service/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	handlers.InitRedis()

	api := r.Group("/api/v1/session")
	{
		api.POST("/", handlers.CreateSession)
		api.GET("/user/:user_id", handlers.GetSessionsByUser)
		api.DELETE("/:session_id", handlers.DeleteSession)
		api.DELETE("/user/:user_id", handlers.DeleteSessionsByUser)
	}
}
