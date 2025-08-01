package api

import (
	_ "github.com/beganov/gingonicserver/internal/api/docs"
	"github.com/beganov/gingonicserver/internal/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RouteRegister(router *gin.Engine) {
	server := NewServer()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.POST("/rooms/", server.createRoom)
	router.GET("/rooms/:id/", server.getRoom)
	router.DELETE("/rooms/:id/", server.deleteRoom) //возможно не надо
	router.PATCH("/rooms/:id/", server.patchRoom)
	router.POST("/rooms/:id/join", server.joinRoom)
	router.DELETE("/rooms/:id/leave", server.leaveRoom)
	router.POST("/rooms/:id/start", server.start)
	router.POST("/rooms/:id/move", server.move)
}

func SetupRouter(zllogger zerolog.Logger) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	router.Use(logger.LoggerMiddleware(zllogger))

	RouteRegister(router)

	return router
}
