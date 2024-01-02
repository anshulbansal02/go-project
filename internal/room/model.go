package room

type Room struct {
	ID           string   `json:"id"`
	Code         string   `json:"code"`
	Type         string   `json:"type"`
	Participants []string `json:"participantIds"`
	Admin        *string  `json:"adminId"`
}

type EventType int

const (
	JoinRequestedEvent EventType = iota
	JoinRequestRejectedEvent
	JoinRequestCancelledEvent
	UserJoinedEvent
	UserLeftEvent
)

type UserEvent struct {
	Type   EventType
	UserId string
	RoomId string
}
