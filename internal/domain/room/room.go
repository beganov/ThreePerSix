package room

import (
	"context"
	"sync"

	"github.com/beganov/gingonicserver/internal/domain/core/game"
	lobbyerror "github.com/beganov/gingonicserver/internal/errors"
)

const maxPlayer = 6

type Room struct { // Room - игровая комната, по сути хранилище для игры
	sync.RWMutex   `json:"-" swaggerignore:"true"`
	Id             int            `json:"id"`             // Уникальный идентификатор комнаты
	MaxPlayerCount int            `json:"maxPlayerCount"` // Максимальное количество игроков
	HostId         int            `json:"hostId"`         // ID игрока, хостящего комнату
	NextPlayerId   int            `json:"nextPlayerId"`   // ID, который будет присвоен следующему игроку
	IsStart        bool           `json:"isStart"`        // Флаг, указывающий, началась ли игра
	Players        map[int]int    `json:"players"`        // Мапа игроков в комнате: [playerID] => playerID
	GameStates     game.GameState `json:"gamestates"`     // Объект игры
}

// RoomUpdate - для изменения максимума игроков патчем
type RoomUpdate struct {
	MaxPlayerCount *int `json:"maxPlayerCount,omitempty"` // Максимальное количество игроков
}

func NewRoom(id int) *Room { // NewRoom создает новую комнату с заданным ID.
	r := &Room{}
	r.Id = id
	r.MaxPlayerCount = maxPlayer
	r.Players = make(map[int]int)
	r.NextPlayerId = 1
	r.HostId = r.NextPlayerId
	r.Players[r.NextPlayerId] = r.NextPlayerId
	r.IsStart = false
	return r
}

func (r *Room) PatchRoom(update RoomUpdate) error { // PatchRoom обновляет параметры комнаты.
	r.Lock()
	defer r.Unlock()
	if r.IsStart { // Ошибка, если игра уже начата
		return lobbyerror.ErrGameAlreadyStarted
	}
	if update.MaxPlayerCount != nil { // Если патч не пустой
		MaxPlayerCount := *update.MaxPlayerCount
		if MaxPlayerCount > 1 && MaxPlayerCount <= 6 && MaxPlayerCount >= len(r.Players) {
			r.MaxPlayerCount = MaxPlayerCount
		} else {
			return lobbyerror.ErrInvalidMaxPlayersCount // Ошибка, если хотим сделать игроков меньше 2, больше 6 или меньше нынешнего числа игроков
		}
	}
	return nil
}

func (r *Room) JoinRoom() (int, error) { // JoinRoom добавляет нового игрока в комнату
	r.Lock()
	defer r.Unlock()
	if r.IsStart { // Ошибка, если игра уже начата
		return 0, lobbyerror.ErrGameAlreadyStarted
	}
	if r.MaxPlayerCount > len(r.Players) {
		r.NextPlayerId++
		r.Players[r.NextPlayerId] = r.NextPlayerId
		return r.NextPlayerId, nil
	}
	return 0, lobbyerror.ErrRoomIsFull
}

func (r *Room) LenRoom() int { // LenRoom возвращает текущее количество игроков в комнате
	r.RLock()
	defer r.RUnlock()
	return len(r.Players)
}

func (r *Room) LeaveRoom(playerId int) error { // LeaveRoom удаляет игрока из комнаты
	r.Lock()
	defer r.Unlock()
	_, isExist := r.Players[playerId]
	if !isExist {
		return lobbyerror.ErrInvalidPlayerID
	}
	if r.IsStart {
		r.GameStates.LeaveGame(playerId)
	} else {
		if playerId == r.HostId { // Если выходит хост
			for i := range r.Players {
				if i != r.HostId {
					r.HostId = i // Делаем хостом кого-то другого из комнаты
					break
				}
			}
		}

	}
	delete(r.Players, playerId)
	return nil
}

func (r *Room) Start(ctx context.Context) (*Room, error) { // Start запускает игру
	r.Lock()
	defer r.Unlock()
	if r.IsStart { // Ошибка, если игра уже начата
		return nil, lobbyerror.ErrGameAlreadyStarted
	}
	r.IsStart = true
	r.GameStates = *r.GameStates.StartGame(r.MaxPlayerCount, r.Players, r, ctx)
	return r, nil
}

func (r *Room) Move(playerId int, playerMove int) error { // Move обрабатывает ход игрока
	r.Lock()
	defer r.Unlock()
	if !r.IsStart { // Ошибка, если игра еще не начата
		return lobbyerror.ErrGameNotStarted
	}
	_, isExist := r.Players[playerId]
	if !isExist {
		return lobbyerror.ErrInvalidPlayerID
	}
	r.GameStates = *r.GameStates.Move(playerId, playerMove)
	return nil
}

func (r *Room) OnGameEnd() { // OnGameEnd вызывается при завершении игры
	r.IsStart = false
	//r.GameStates = game.GameState{}
}
