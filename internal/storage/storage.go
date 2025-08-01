package storage

import (
	"context"
	"sync"

	"github.com/beganov/gingonicserver/internal/domain/room"
	lobbyerror "github.com/beganov/gingonicserver/internal/errors"
)

type Storage struct { // Storage — хранилище игровых комнат
	sync.RWMutex
	rooms      map[int]*room.Room // мапа всех активных комнат
	nextRoomId int                // следующий ID для новой комнаты
}

func NewStorage() *Storage { // NewStorage создаёт новое хранилище комнат
	s := &Storage{}
	s.rooms = make(map[int]*room.Room)
	s.nextRoomId = 1
	return s
}

func (s *Storage) CreateRoom() (int, int) { // CreateRoom создаёт новую игровую комнату и возвращает её ID и ID хоста
	s.Lock()
	defer s.Unlock()

	room := room.NewRoom(s.nextRoomId)
	s.rooms[s.nextRoomId] = room
	s.nextRoomId++
	if s.nextRoomId == 1<<32 { //лениво (Если доходим до какого то числа комнат, то сбрасываем до 0)
		//(видимо игнорируем, что при этом будет, если такая комната уже есть)
		s.nextRoomId = 0
	}
	return room.Id, room.HostId
}

func (s *Storage) GetRoom(roomId int) (*room.Room, error) { // GetRoom возвращает комнату по ID, если она существует
	s.RLock() //Чет с локами надо разобраться будет
	room, isExist := s.rooms[roomId]
	s.RUnlock()
	if !isExist {
		return nil, lobbyerror.ErrInvalidRoomID
	}
	return room, nil
}

func (s *Storage) DeleteRoom(roomId int) error { // DeleteRoom удаляет комнату по ID
	s.Lock()
	defer s.Unlock()

	_, isExist := s.rooms[roomId]
	if !isExist {
		return lobbyerror.ErrInvalidRoomID
	}
	delete(s.rooms, roomId)
	return nil
}

func (s *Storage) PatchRoom(roomId int, update room.RoomUpdate) error { // PatchRoom обновляет состояние комнаты с заданным ID с помощью update
	s.RLock()
	room, isExist := s.rooms[roomId]
	s.RUnlock()
	if !isExist {
		return lobbyerror.ErrInvalidRoomID
	}
	return room.PatchRoom(update)
}

func (s *Storage) JoinRoom(roomId int) (int, error) { // JoinRoom добавляет игрока в комнату и возвращает его ID
	s.RLock()
	room, isExist := s.rooms[roomId]
	s.RUnlock()
	if !isExist {
		return 0, lobbyerror.ErrInvalidRoomID
	}
	return room.JoinRoom()
}

func (s *Storage) LeaveRoom(roomId int, playerId int) error { // LeaveRoom удаляет игрока из комнаты; удаляет комнату, если она пуста
	s.Lock()
	defer s.Unlock()
	room, isExist := s.rooms[roomId]
	if !isExist {
		return lobbyerror.ErrInvalidRoomID
	}
	err := room.LeaveRoom(playerId)
	if room.LenRoom() == 0 { //если комната пуста - удаляем ее
		delete(s.rooms, roomId)
	}
	return err
}

func (s *Storage) Start(roomId int, ctx context.Context) (*room.Room, error) { // Start запускает игру в комнате с заданным ID
	s.RLock()
	room, isExist := s.rooms[roomId]
	s.RUnlock()
	if !isExist {
		return nil, lobbyerror.ErrInvalidRoomID
	}
	return room.Start(ctx)
}

func (s *Storage) Move(roomId int, playerId int, playerMove int) error { // Move выполняет ход игрока в заданной комнате
	s.RLock()
	room, isExist := s.rooms[roomId]
	s.RUnlock()
	if !isExist {
		return lobbyerror.ErrInvalidRoomID
	}
	return room.Move(playerId, playerMove)
}
