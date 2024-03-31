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

func NewChatEventsExchange(chatService *chat.ChatService, roomService *room.RoomService, wsManager *websockets.WebSocketManager, clientMap *ClientMap) *ChatEventsExchange {
	return &ChatEventsExchange{
		chatService: chatService,
		roomService: roomService,
		wsManager:   wsManager,
		clientMap:   clientMap,
	}
}

func (e *ChatEventsExchange) Listen() {

	go func() {
		for {
			msg := <-e.chatService.ChatMessageChannel

			e.handleIE_ChatMessage(msg)
		}
	}()

	e.wsManager.AddObserver(events.Chat.ChatMessage, e.handleCE_ChatMessage)

}

func (e *ChatEventsExchange) handleIE_ChatMessage(msg *chat.ChatMessage) {
	userIds, err := e.roomService.GetRoomUsers(context.Background(), msg.ConversationId)
	if err != nil {
		return
	}

	clientIds := e.clientMap.GetClientIds(userIds)
	e.wsManager.Multicast(clientIds, nil, websockets.NewNotification(events.Chat.ChatMessage, events.OutgoingChatMessageData{
		ID:             msg.ID,
		Content:        msg.Content,
		Meta:           msg.Meta,
		UserId:         msg.UserId,
		Timestamp:      msg.Timestamp,
		ConversationId: msg.ConversationId,
	}))
}

func (e *ChatEventsExchange) handleCE_ChatMessage(message websockets.IncomingWebSocketMessage, client *websockets.Client) {

	data := events.IncomingChatMessageData{}
	if err := message.Payload.Assert(&data); err != nil {
		return
	}

	ctx := context.Background()

	userId, exists := e.clientMap.GetUserId(client.ID)
	if !exists {
		return
	}

	roomId, err := e.roomService.GetUserRoomId(ctx, userId)
	if err != nil {
		return
	}

	e.chatService.CreateMessage(ctx, data.Content, data.Meta, userId, roomId)
}
