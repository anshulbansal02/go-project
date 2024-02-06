package events

import "anshulbansal02/scribbly/pkg/websockets"

var Room = struct {
	JoinRequest       websockets.Event
	CancelJoinRequest websockets.Event
	UserJoined        websockets.Event
	UserLeft          websockets.Event
}{
	JoinRequest: "join_request",
	UserJoined:  "user_joined",
	UserLeft:    "user_left",
}

type RequestData struct {
	Type   string `msgpack:"type"`
	UserId string `msgpack:"userId"`
	RoomId string `msgpack:"roomId"`
}

type RoomUserData struct {
	RoomId string `msgpack:"roomId"`
	UserId string `msgpack:"userId"`
}
