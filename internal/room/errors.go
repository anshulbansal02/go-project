package room

import "errors"

var (
	ErrRoomNotFound      = errors.New("room not found")
	ErrUserAlreadyInRoom = errors.New("user is already in another room")
)
