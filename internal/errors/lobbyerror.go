package lobbyerror

import "errors"

var (
	ErrInvalidRoomID          = errors.New("invalid room ID")
	ErrInvalidPlayerID        = errors.New("invalid player ID")
	ErrGameAlreadyStarted     = errors.New("game has already started")
	ErrGameNotStarted         = errors.New("game has not started yet")
	ErrRoomIsFull             = errors.New("room is full")
	ErrInvalidMaxPlayersCount = errors.New("invalid max players count")
)
