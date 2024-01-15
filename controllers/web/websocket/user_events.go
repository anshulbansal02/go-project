package exchange

import (
	"anshulbansal02/scribbly/events"
	"anshulbansal02/scribbly/internal/user"
	"anshulbansal02/scribbly/pkg/websockets"
	"fmt"
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

	e.wsManager.AddObserver(events.User.AssociateClient, func(m websockets.IncomingWebSocketMessage, c *websockets.Client) {

		k := events.AssociateClientData{}
		err := m.Payload.Assert(&k)

		if err != nil {
			fmt.Println("Cannot type assert associate client data")
			return
		}

		fmt.Println(c.ID, k)

		userId, err := e.userService.VerifyUserToken(k.UserSecret)
		if err == nil {
			e.clientMap.Add(c.ID, userId)
		}
	})

}
