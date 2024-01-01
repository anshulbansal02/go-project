package exchange

import (
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
	handleUserRoomEvents := func() {
		for {
			event := <-e.roomService.UserEventsChannel

			switch event.Type {
			case room.JoinRequestEvent:
				go e.handleIE_JoinRequest(event.UserId, event.RoomId)
			case room.CancelJoinRequestEvent:
				go e.handleIE_CancelRequest(event.UserId, event.RoomId)
			case room.UserJoinEvent:
				go e.handleIE_RoomJoined(event.UserId, event.RoomId)
			case room.UserLeaveEvent:
				go e.handleIE_RoomLeft(event.UserId, event.RoomId)
			}

		}

	}

	// Handle Events from clients
	e.wsManager.OnEvent("cancelRequest", e.handleCE_CancelRequest)

	go handleUserRoomEvents()

}

func (e *RoomEventsExchange) handleIE_JoinRequest(userId, roomId string) {
	adminId := e.roomService.GetRoomAdmin(context.Background(), roomId)

	if adminId != nil {
		e.wsManager.Send(roomId, "joinRequest:{userId}")
	}
}

func (e *RoomEventsExchange) handleIE_CancelRequest(userId, roomId string) {
	adminId := e.roomService.GetRoomAdmin(context.Background(), roomId)

	if adminId != nil {
		e.wsManager.Send(roomId, "cancelRequest:{userId}")
	}
}

func (e *RoomEventsExchange) handleIE_RoomLeft(userId, roomId string) {
	userIds, err := e.roomService.GetRoomUsers(context.Background(), roomId)
	if err != nil {
		return
	}

	e.wsManager.Multicast(userIds, nil, "userLEft:{userId}")

}

func (e *RoomEventsExchange) handleIE_RoomJoined(userId, roomId string) {
	userIds, err := e.roomService.GetRoomUsers(context.Background(), roomId)
	if err != nil {
		return
	}

	e.wsManager.Multicast(userIds, nil, "userJoined:{userId}")

}

func (e *RoomEventsExchange) handleCE_CancelRequest(client *websockets.Client, event websockets.WebSocketEvent) {
	userId, ok := e.clientMap.GetUserId(client.ID)
	if !ok {
		return
	}

	e.roomService.CancelJoinRequest(context.Background(), userId)
}
