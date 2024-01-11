package exchange

import (
	"anshulbansal02/scribbly/events"
	"anshulbansal02/scribbly/internal/user"
	"anshulbansal02/scribbly/pkg/websockets"
	"encoding/json"
)

type UserEventsExchange struct {
	userService *user.UserService
	wsManager   *websockets.WebSocketManager
	clientMap   *ClientMap
}

func NewUserEventsExchange(userService *user.UserService, wsManager *websockets.WebSocketManager) *UserEventsExchange {
	return &UserEventsExchange{
		userService: userService,
		wsManager:   wsManager,
		clientMap:   NewClientMap(),
	}
}

func (e *UserEventsExchange) Listen() {

	e.wsManager.AddObserver(events.User.AssociateClient, func(m websockets.WebSocketMessage, c *websockets.Client) {
		p := events.AssociateClientData{}
		if err := json.Unmarshal(m.Payload.([]byte), &p); err != nil {
			return
		}

		userId, err := e.userService.VerifyUserToken(p.UserSecret)
		if err == nil {
			e.clientMap.Add(c.ID, userId)
		}
	})

}
