package main

import (
	"fmt"

	"github.com/beganov/gingonicserver/internal/api"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	api.RouteRegister(router)
	if err := router.Run("localhost:8080"); err != nil {
		fmt.Println("Failed to run server:", err)
	}
}
