package api

import (
	"net/http"
	"strconv"

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

// createRoom создает новую комнату и возвращает её ID и ID создателя
// @Summary      Создать комнату
// @Description  Создает новую игровую комнату и возвращает ID комнаты и ID игрока-организатора
// @Tags         rooms
// @Produce      json
// @Success      201  {object}  map[string]string  "roomId и playerId"
// @Router       /rooms/ [post]
func (gs *gameServer) createRoom(c *gin.Context) {
	roomId, playerID := gs.store.CreateRoom()
	c.JSON(http.StatusCreated, gin.H{"roomId": roomId, "playerId": playerID})
}

// getRoom возвращает данные комнаты по ID
// @Summary      Получить информацию о комнате
// @Description  Возвращает информацию о комнате по её ID
// @Tags         rooms
// @Produce      json
// @Param        id   path      int  true  "ID комнаты"
// @Success      200  {object}  room.Room  "Данные комнаты"
// @Failure      400  {string}  string     "Некорректный ID"
// @Failure      404  {string}  string     "Комната не найдена"
// @Router       /rooms/{id}/ [get]
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

// joinRoom добавляет игрока в комнату и возвращает его ID
// @Summary      Присоединиться к комнате
// @Description  Позволяет игроку присоединиться к комнате по ID и возвращает его playerId
// @Tags         rooms
// @Produce      json
// @Param        id   path      int  true  "ID комнаты"
// @Success      200  {object}  map[string]string  "playerId"
// @Failure      400  {string}  string             "Некорректный ID комнаты"
// @Failure      409  {string}  string             "Ошибка присоединения (например, комната полна)"
// @Router       /rooms/{id}/join [post]swag init
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
	err = gs.store.Move(roomId, currentPlayer.Id, currentPlayer.Move)
	if err != nil {
		c.String(http.StatusConflict, formatError(err))
		return
	}
	//debugPrintGameState(game)
	//time.Sleep(100) // эта вещь точно работала, строчку ниже не тестил, но по логике должна
	c.JSON(http.StatusOK, gin.H{"message": "move accepted"})
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
