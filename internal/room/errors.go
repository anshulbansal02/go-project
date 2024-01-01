package room

import "errors"

var (
	ErrRoomNotFound      = errors.New("Room not found")
	ErrUserAlreadyInRoom = errors.New("User is already in another room")
)
