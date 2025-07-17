package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/beganov/gingonicserver/internal/player"
	"github.com/beganov/gingonicserver/internal/room"
	"github.com/beganov/gingonicserver/internal/storage"
	"github.com/gin-gonic/gin"
)

type gameServer struct {
	store *storage.Storage
}

func NewServer() *gameServer {
	store := storage.NewStorage()
	return &gameServer{store: store}
}

func (gs *gameServer) createRoom(c *gin.Context) {
	roomId, playerID := gs.store.CreateRoom()
	c.JSON(http.StatusCreated, gin.H{"roomId": roomId, "playerId": playerID})
}

func (gs *gameServer) getRoom(c *gin.Context) {
	roomId, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	game, err := gs.store.GetRoom(roomId)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, game)
}

func (gs *gameServer) deleteRoom(c *gin.Context) {
	roomId, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	err = gs.store.DeleteRoom(roomId)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

func (gs *gameServer) patchRoom(c *gin.Context) {
	roomId, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	var roomPatch room.RoomUpdate
	if err := c.ShouldBindJSON(&roomPatch); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	err = gs.store.PatchRoom(roomId, roomPatch)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

func (gs *gameServer) joinRoom(c *gin.Context) {
	roomId, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	playerID, err := gs.store.JoinRoom(roomId)
	if err != nil {
		c.String(http.StatusConflict, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"playerId": playerID})
}

func (gs *gameServer) leaveRoom(c *gin.Context) {
	roomId, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	var currentPlayer player.Player
	if err := c.ShouldBindJSON(&currentPlayer); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	err = gs.store.LeaveRoom(roomId, currentPlayer.Id)
	if err != nil {
		c.String(http.StatusConflict, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

func (gs *gameServer) start(c *gin.Context) {
	roomId, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	game, err := gs.store.Start(roomId)
	if err != nil {
		c.String(http.StatusConflict, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"game": game})
}

func (gs *gameServer) move(c *gin.Context) {
	roomId, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	var currentPlayer player.Player
	if err := c.ShouldBindJSON(&currentPlayer); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	game, err := gs.store.Move(roomId, currentPlayer.Id, currentPlayer.Move)
	if err != nil {
		c.String(http.StatusConflict, err.Error())
		return
	}
	//debugPrintGameState(game)
	time.Sleep(100)
	c.JSON(http.StatusOK, gin.H{"game": game})
}

func main() {
	router := gin.Default()
	server := NewServer()

	router.POST("/rooms/", server.createRoom)
	router.GET("/rooms/:id/", server.getRoom)
	router.DELETE("/rooms/:id/", server.deleteRoom) //возможно не надо
	router.PATCH("/rooms/:id/", server.patchRoom)
	router.POST("/rooms/:id/join", server.joinRoom)
	router.DELETE("/rooms/:id/leave", server.leaveRoom)
	router.POST("/rooms/:id/start", server.start)
	router.POST("/rooms/:id/move", server.move)
	router.Run("localhost:8080")
}
