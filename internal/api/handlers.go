package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/beganov/gingonicserver/internal/domain/player"
	"github.com/beganov/gingonicserver/internal/domain/room"
	"github.com/beganov/gingonicserver/internal/storage"
	"github.com/gin-gonic/gin"

	lobbyerror "github.com/beganov/gingonicserver/internal/errors"
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
		c.String(http.StatusBadRequest, formatError(lobbyerror.ErrInvalidRoomID))
		return
	}

	room, err := gs.store.GetRoom(roomId)
	if err != nil {
		c.String(http.StatusNotFound, formatError(err))
		return
	}

	c.JSON(http.StatusOK, room)
}

func (gs *gameServer) deleteRoom(c *gin.Context) {
	roomId, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.String(http.StatusBadRequest, formatError(lobbyerror.ErrInvalidRoomID))
		return
	}

	err = gs.store.DeleteRoom(roomId)
	if err != nil {
		c.String(http.StatusNotFound, formatError(err))
		return
	}

	c.Status(http.StatusNoContent)
}

func (gs *gameServer) patchRoom(c *gin.Context) {
	roomId, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.String(http.StatusBadRequest, formatError(lobbyerror.ErrInvalidRoomID))
		return
	}
	var roomPatch room.RoomUpdate
	if err := c.ShouldBindJSON(&roomPatch); err != nil {
		c.String(http.StatusBadRequest, formatError(lobbyerror.ErrInvalidMaxPlayersCount))
		return
	}

	err = gs.store.PatchRoom(roomId, roomPatch)
	if err != nil {
		c.String(http.StatusNotFound, formatError(err))
		return
	}

	c.Status(http.StatusNoContent)
}

func (gs *gameServer) joinRoom(c *gin.Context) {
	roomId, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.String(http.StatusBadRequest, formatError(lobbyerror.ErrInvalidRoomID))
		return
	}
	playerID, err := gs.store.JoinRoom(roomId)
	if err != nil {
		c.String(http.StatusConflict, formatError(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"playerId": playerID})
}

func (gs *gameServer) leaveRoom(c *gin.Context) {
	roomId, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.String(http.StatusBadRequest, formatError(lobbyerror.ErrInvalidRoomID))
		return
	}
	var currentPlayer player.Player
	if err := c.ShouldBindJSON(&currentPlayer); err != nil {
		c.String(http.StatusBadRequest, formatError(lobbyerror.ErrInvalidPlayerID))
		return
	}
	err = gs.store.LeaveRoom(roomId, currentPlayer.Id)
	if err != nil {
		c.String(http.StatusConflict, formatError(err))
		return
	}
	c.Status(http.StatusNoContent)
}

func (gs *gameServer) start(c *gin.Context) {
	roomId, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.String(http.StatusBadRequest, formatError(lobbyerror.ErrInvalidRoomID))
		return
	}
	room, err := gs.store.Start(roomId)
	if err != nil {
		c.String(http.StatusConflict, formatError(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"room": room})
}

func (gs *gameServer) move(c *gin.Context) {
	roomId, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.String(http.StatusBadRequest, formatError(lobbyerror.ErrInvalidRoomID))
		return
	}
	var currentPlayer player.Player
	if err := c.ShouldBindJSON(&currentPlayer); err != nil {
		c.String(http.StatusBadRequest, formatError(lobbyerror.ErrInvalidPlayerID))
		return
	}
	room, err := gs.store.Move(roomId, currentPlayer.Id, currentPlayer.Move)
	if err != nil {
		c.String(http.StatusConflict, formatError(err))
		return
	}
	//debugPrintGameState(game)
	time.Sleep(100)
	c.JSON(http.StatusOK, gin.H{"room": room})
}

func formatError(err error) string {
	switch err {
	case lobbyerror.ErrRoomIsFull:
		return "Извините, комната уже полная."
	case lobbyerror.ErrGameAlreadyStarted:
		return "Игра уже началась. Присоединение невозможно."
	case lobbyerror.ErrGameNotStarted:
		return "Игра ещё не началась. Подождите, пока организатор начнёт игру."
	case lobbyerror.ErrInvalidRoomID:
		return "Некорректный ID комнаты. Проверьте правильность ввода."
	case lobbyerror.ErrInvalidPlayerID:
		return "Некорректный ID игрока. Попробуйте перезайти или обновить страницу."
	case lobbyerror.ErrInvalidMaxPlayersCount:
		return "Недопустимое максимальное число для игроков в комнате."
	default:
		return "Внутренняя ошибка сервера. Попробуйте позже."
	}
}
