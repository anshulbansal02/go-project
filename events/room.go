package events

import "anshulbansal02/scribbly/pkg/websockets"

var Room = struct {
	JoinRequest       websockets.Event
	CancelJoinRequest websockets.Event
	UserJoined        websockets.Event
	UserLeft          websockets.Event
}{
	JoinRequest:       "join_request",
	CancelJoinRequest: "cancel_join_request",
	UserJoined:        "user_joined",
	UserLeft:          "user_left",
}

type RequestData struct {
	UserId string
}

type RoomUserData struct {
	RoomId string
	UserId string
}
