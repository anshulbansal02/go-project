package exchange

import (
	"anshulbansal02/scribbly/events"
	"anshulbansal02/scribbly/internal/chat"
	"anshulbansal02/scribbly/internal/room"
	"anshulbansal02/scribbly/pkg/websockets"

	"context"
)

type ChatEventsExchange struct {
	chatService *chat.ChatService
	roomService *room.RoomService
	wsManager   *websockets.WebSocketManager
	clientMap   *ClientMap
}

func NewChatEventsExchange(chatService *chat.ChatService, wsManager *websockets.WebSocketManager, clientMap *ClientMap) *ChatEventsExchange {
	return &ChatEventsExchange{
		chatService: chatService,
		wsManager:   wsManager,
		clientMap:   clientMap,
	}
}

func (e *ChatEventsExchange) Listen() {

	go func() {
		for {
			msg := <-e.chatService.ChatMessageChannel

			userIds, err := e.roomService.GetRoomUsers(context.Background(), msg.ConversationId)
			if err != nil {
				return
			}

			clientIds := e.clientMap.GetClientIds(userIds)

			e.wsManager.Multicast(clientIds, nil, websockets.NewNotification(events.Chat.ChatMessage, msg))

		}
	}()

}
