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

func NewRoomEventsExchange(roomService *room.RoomService, wsManager *websockets.WebSocketManager) *RoomEventsExchange {
	return &RoomEventsExchange{
		roomService: roomService,
		wsManager:   wsManager,
		clientMap:   NewClientMap(),
	}
}

func (e *RoomEventsExchange) Listen() {

	// Handle Events from Internal channel
	handleInternalUserRoomEvents := func() {
		for {
			event := <-e.roomService.UserEventsChannel

			switch event.Type {
			case room.JoinRequestedEvent:
				go e.handleIE_JoinRequest(event.UserId, event.RoomId)
			case room.JoinRequestCancelledEvent:
				go e.handleIE_CancelRequest(event.UserId, event.RoomId)
			case room.UserJoinedEvent:
				go e.handleIE_UserJoined(event.UserId, event.RoomId)
			case room.UserLeftEvent:
				go e.handleIE_UserLeft(event.UserId, event.RoomId)
			}

		}

	}

	// Handle Events from clients
	e.wsManager.AddObserver(events.Room.CancelJoinRequest, e.handleCE_CancelRequest)

	go handleInternalUserRoomEvents()

}

// System Event Handlers

func (e *RoomEventsExchange) handleIE_JoinRequest(userId, roomId string) {
	adminId := e.roomService.GetRoomAdmin(context.Background(), roomId)

	if adminId != nil {
		clientId, ok := e.clientMap.GetClientId(*adminId)
		if ok {
			e.wsManager.EmitTo(clientId, websockets.NewNotification(events.Room.JoinRequest, events.RequestData{
				UserId: userId,
			}))
		}
	}
}

func (e *RoomEventsExchange) handleIE_CancelRequest(userId, roomId string) {
	adminId := e.roomService.GetRoomAdmin(context.Background(), roomId)

	if adminId != nil {
		clientId, ok := e.clientMap.GetClientId(*adminId)
		if ok {
			e.wsManager.EmitTo(clientId, websockets.NewNotification(events.Room.CancelJoinRequest, events.RequestData{
				UserId: userId,
			}))
		}
	}
}

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

func (e *RoomEventsExchange) handleIE_UserJoined(userId, roomId string) {
	userIds, err := e.roomService.GetRoomUsers(context.Background(), roomId)
	if err != nil {
		return
	}

	clientIds := e.clientMap.GetClientIds(userIds)

	e.wsManager.Multicast(clientIds, nil, websockets.NewNotification(events.Room.UserJoined, events.RoomUserData{
		UserId: userId,
	}))
}

// Client Event Handlers

func (e *RoomEventsExchange) handleCE_CancelRequest(message websockets.IncomingWebSocketMessage, client *websockets.Client) {
	userId, ok := e.clientMap.GetUserId(client.ID)
	if !ok {
		return
	}

	e.roomService.CancelJoinRequest(context.Background(), userId)
}
