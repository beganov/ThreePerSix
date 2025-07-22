package storage

import (
	"sync"

	gameerror "github.com/beganov/gingonicserver/internal/errors"
	"github.com/beganov/gingonicserver/internal/room"
)

type Storage struct {
	sync.RWMutex
	rooms      map[int]*room.Room
	nextRoomId int
}

func NewStorage() *Storage {
	s := &Storage{}
	s.rooms = make(map[int]*room.Room)
	s.nextRoomId = 1
	return s
}

func (s *Storage) CreateRoom() (int, int) {
	s.Lock()
	defer s.Unlock()

	room := room.NewRoom(s.nextRoomId)
	s.rooms[s.nextRoomId] = room
	s.nextRoomId++
	if s.nextRoomId == 1<<32 { //лениво
		s.nextRoomId = 0
	}
	return room.Id, room.HostId
}

func (s *Storage) GetRoom(roomId int) (*room.Room, error) {
	s.RLock() //Чет с локами надо разобраться будет
	room, isExist := s.rooms[roomId]
	s.RUnlock()
	if !isExist {
		return nil, gameerror.ErrIncorrectRoomId
	}
	return room, nil
}

func (s *Storage) DeleteRoom(roomId int) error {
	s.Lock()
	defer s.Unlock()

	_, isExist := s.rooms[roomId]
	if !isExist {
		return gameerror.ErrIncorrectRoomId
	}
	delete(s.rooms, roomId)
	return nil
}

func (s *Storage) PatchRoom(roomId int, update room.RoomUpdate) error {
	s.RLock()
	room, isExist := s.rooms[roomId]
	s.RUnlock()
	if !isExist {
		return gameerror.ErrIncorrectRoomId
	}
	return room.PatchRoom(update)
}

func (s *Storage) JoinRoom(roomId int) (int, error) {
	s.RLock()
	room, isExist := s.rooms[roomId]
	s.RUnlock()
	if !isExist {
		return 0, gameerror.ErrIncorrectRoomId
	}
	return room.JoinRoom()
}

func (s *Storage) LeaveRoom(roomId int, playerId int) error {
	s.Lock()
	defer s.Unlock()
	room, isExist := s.rooms[roomId]
	if !isExist {
		return gameerror.ErrIncorrectRoomId
	}
	err := room.LeaveRoom(playerId)
	if room.LenRoom() == 0 {
		delete(s.rooms, roomId)
	}
	return err
}

func (s *Storage) Start(roomId int) (*room.Room, error) {
	s.RLock()
	room, isExist := s.rooms[roomId]
	s.RUnlock()
	if !isExist {
		return nil, gameerror.ErrIncorrectRoomId
	}
	return room.Start()
}

func (s *Storage) Move(roomId int, playerId int, playerMove int) (*room.Room, error) { //надо вернуть GameState
	s.RLock()
	room, isExist := s.rooms[roomId]
	s.RUnlock()
	if !isExist {
		return nil, gameerror.ErrIncorrectRoomId
	}
	return room.Move(playerId, playerMove)
}
