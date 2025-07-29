package room

import (
	"sync"

	"github.com/beganov/gingonicserver/internal/domain/core/game"
	lobbyerror "github.com/beganov/gingonicserver/internal/errors"
)

const maxPlayer = 6

type Room struct {
	sync.RWMutex
	Id             int            `json:"id"`
	MaxPlayerCount int            `json:"maxPlayerCount"`
	HostId         int            `json:"hostId"`
	NextPlayerId   int            `json:"nextPlayerId"`
	IsStart        bool           `json:"isStart"`
	Players        map[int]int    `json:"players"`
	GameStates     game.GameState `json:"gamestates"`
}

type RoomUpdate struct {
	MaxPlayerCount *int `json:"maxPlayerCount,omitempty"`
}

func NewRoom(id int) *Room {
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

func (r *Room) PatchRoom(update RoomUpdate) error {
	r.Lock()
	defer r.Unlock()
	if r.IsStart {
		return lobbyerror.ErrStart
	}
	if update.MaxPlayerCount != nil {
		MaxPlayerCount := *update.MaxPlayerCount
		if MaxPlayerCount > 1 && MaxPlayerCount <= 6 && MaxPlayerCount >= len(r.Players) {
			r.MaxPlayerCount = MaxPlayerCount
		} else {
			return lobbyerror.ErrMaxPlayerCount
		}
	}
	return nil
}

func (r *Room) JoinRoom() (int, error) {
	r.Lock()
	defer r.Unlock()
	if r.IsStart {
		return 0, lobbyerror.ErrStart
	}
	if r.MaxPlayerCount > len(r.Players) {
		r.NextPlayerId++
		r.Players[r.NextPlayerId] = r.NextPlayerId
		return r.NextPlayerId, nil
	}
	return 0, lobbyerror.ErrFullRoom
}

func (r *Room) LenRoom() int {
	r.RLock()
	defer r.RUnlock()
	return len(r.Players)
}

func (r *Room) LeaveRoom(playerId int) error {
	r.Lock()
	defer r.Unlock()
	_, isExist := r.Players[playerId]
	if !isExist {
		return lobbyerror.ErrIncorrectRoomId
	}
	if r.IsStart {
		r.GameStates.LeaveGame(playerId)
	} else {
		if playerId == r.HostId {
			for i := range r.Players {
				if i != r.HostId {
					r.HostId = i
					break
				}
			}
		}

	}
	delete(r.Players, playerId)
	return nil
}

func (r *Room) Start() (*Room, error) {
	r.Lock()
	defer r.Unlock()
	if r.IsStart {
		return nil, lobbyerror.ErrStart
	}
	r.IsStart = true
	r.GameStates = *r.GameStates.StartGame(r.MaxPlayerCount, r.Players, r)
	return r, nil
}

func (r *Room) Move(playerId int, playerMove int) (*Room, error) { //надо вернуть GameState
	r.Lock()
	defer r.Unlock()
	if !r.IsStart {
		return nil, lobbyerror.ErrNotStart
	}
	_, isExist := r.Players[playerId]
	if !isExist {
		return nil, lobbyerror.ErrIncorrectPlayerId
	}
	r.GameStates = *r.GameStates.Move(playerId, playerMove)
	return r, nil
}

func (r *Room) OnGameEnd() {
	r.IsStart = false

	//r.GameStates = game.GameState{}
}
