package main

import (
	"session-service/config"
	"session-service/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()

	r := gin.Default()
	routes.SetupRoutes(r)

	port := config.GetPort()
	r.Run(":" + port)
}
