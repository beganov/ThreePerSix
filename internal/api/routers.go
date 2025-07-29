package api

import "github.com/gin-gonic/gin"

func RouteRegister(router *gin.Engine) {
	server := NewServer()

	router.POST("/rooms/", server.createRoom)
	router.GET("/rooms/:id/", server.getRoom)
	router.DELETE("/rooms/:id/", server.deleteRoom) //возможно не надо
	router.PATCH("/rooms/:id/", server.patchRoom)
	router.POST("/rooms/:id/join", server.joinRoom)
	router.DELETE("/rooms/:id/leave", server.leaveRoom)
	router.POST("/rooms/:id/start", server.start)
	router.POST("/rooms/:id/move", server.move)
}
