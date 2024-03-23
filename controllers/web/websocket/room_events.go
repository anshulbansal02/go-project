package exchange

import (
	"anshulbansal02/scribbly/events"
	"anshulbansal02/scribbly/internal/room"
	"anshulbansal02/scribbly/pkg/websockets"

	"context"
)

type RoomEventsExchange struct {
	roomService *room.RoomService
	wsManager   *websockets.WebSocketManager
	clientMap   *ClientMap
}

func NewRoomEventsExchange(roomService *room.RoomService, wsManager *websockets.WebSocketManager, clientMap *ClientMap) *RoomEventsExchange {
	return &RoomEventsExchange{
		roomService: roomService,
		wsManager:   wsManager,
		clientMap:   clientMap,
	}
}

func (e *RoomEventsExchange) Listen() {
	// Handle Events from Internal channel
	go func() {
		for {
			event := <-e.roomService.UserEventsChannel

			switch event.Type {
			case room.JoinRequestedEvent:
				go e.handleIE_JoinRequest(event.UserId, event.RoomId)
			case room.JoinRequestCancelledEvent:
				go e.handleIE_CancelRequest(event.UserId, event.RoomId)
			case room.JoinRequestRejectedEvent:
				go e.handleIE_RejectRequest(event.UserId, event.RoomId)
			case room.UserJoinedEvent:
				go e.handleIE_UserJoined(event.UserId, event.RoomId)
			case room.UserLeftEvent:
				go e.handleIE_UserLeft(event.UserId, event.RoomId)
			}
		}
	}()

	// Handle Events from Clients
	e.wsManager.AddObserver(events.Room.JoinRequest, e.handleCE_ActionOnRequest)

}

// [Internal Event] - New user Join Request for a room
func (e *RoomEventsExchange) handleIE_JoinRequest(userId, roomId string) {
	adminId := e.roomService.GetRoomAdmin(context.Background(), roomId)

	if adminId != nil {
		clientId, ok := e.clientMap.GetClientId(*adminId)
		if ok {
			e.wsManager.EmitTo(clientId, websockets.NewNotification(events.Room.JoinRequest, events.RequestData{
				Type:   "request",
				UserId: userId,
			}))
		}
	}
}

// [Internal Event] - Cancel user Join Request for a room
func (e *RoomEventsExchange) handleIE_CancelRequest(userId, roomId string) {
	adminId := e.roomService.GetRoomAdmin(context.Background(), roomId)

	if adminId != nil {
		clientId, ok := e.clientMap.GetClientId(*adminId)
		if ok {
			e.wsManager.EmitTo(clientId, websockets.NewNotification(events.Room.CancelJoinRequest, events.RequestData{
				Type:   "cancel",
				UserId: userId,
			}))
		}
	}
}

// [Internal Event] - Reject user Join Request for a room
func (e *RoomEventsExchange) handleIE_RejectRequest(userId, roomId string) {
	clientId, ok := e.clientMap.GetClientId(userId)
	if ok {
		e.wsManager.EmitTo(clientId, websockets.NewNotification(events.Room.CancelJoinRequest, events.RequestData{
			Type: "reject",
		}))
	}
}

// [Internal Event] - User Join Request was accepted and user joined a room
func (e *RoomEventsExchange) handleIE_UserJoined(userId, roomId string) {
	userIds, err := e.roomService.GetRoomUsers(context.Background(), roomId)
	if err != nil {
		return
	}

	clientIds := e.clientMap.GetClientIds(userIds)
	userClientId, _ := e.clientMap.GetClientId(userId)

	e.wsManager.Multicast(clientIds, []string{userClientId}, websockets.NewNotification(events.Room.UserJoined, events.RoomUserData{
		UserId: userId,
	}))

	e.wsManager.EmitTo(userClientId, websockets.NewNotification(events.Room.JoinRequest, events.RequestData{
		Type:   "accept",
		UserId: userId,
		RoomId: roomId,
	}))
}

// [Internal Event] - User left a room
func (e *RoomEventsExchange) handleIE_UserLeft(userId, roomId string) {
	userIds, err := e.roomService.GetRoomUsers(context.Background(), roomId)
	if err != nil {
		return
	}

	clientIds := e.clientMap.GetClientIds(userIds)

	e.wsManager.Multicast(clientIds, []string{userId}, websockets.NewNotification(events.Room.UserLeft, events.RoomUserData{
		UserId: userId,
	}))
}

// [Client Event] - There was an action on Join Request either by Admin or by Requestor
func (e *RoomEventsExchange) handleCE_ActionOnRequest(message websockets.IncomingWebSocketMessage, client *websockets.Client) {
	data := events.RequestData{}
	if err := message.Payload.Assert(&data); err != nil {
		return
	}

	ctx := context.Background()

	// check if requestor is admin of the room
	requestorId, _ := e.clientMap.GetUserId(client.ID)

	if data.Type == "accept" || data.Type == "reject" {
		requestorRoomId, err := e.roomService.GetUserRoomId(ctx, requestorId)
		if err != nil {
			return
		}
		adminId := e.roomService.GetRoomAdmin(ctx, requestorRoomId)
		if adminId == nil || *adminId != requestorId {
			return
		}
	}

	switch data.Type {
	case "accept":
		err := e.roomService.AcceptJoinRequest(ctx, data.UserId)
		if err != nil {
			return
		}
	case "reject":
		err := e.roomService.RejectJoinRequest(ctx, data.UserId)
		if err != nil {
			return
		}
	case "cancel":
		err := e.roomService.CancelJoinRequest(ctx, requestorId)
		if err != nil {
			return
		}
	}

}
