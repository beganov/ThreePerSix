package lobbyerror

import "errors"

var ErrIncorrectRoomId = errors.New("incorrect room id")
var ErrIncorrectPlayerId = errors.New("incorrect player id")
var ErrStart = errors.New("game already started")
var ErrNotStart = errors.New("game not sterted yet")
var ErrFullRoom = errors.New("room already full")
var ErrMaxPlayerCount = errors.New("incorrect Max Player setting")
