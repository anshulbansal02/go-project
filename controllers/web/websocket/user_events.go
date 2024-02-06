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

func NewUserEventsExchange(userService *user.UserService, wsManager *websockets.WebSocketManager, clientMap *ClientMap) *UserEventsExchange {
	return &UserEventsExchange{
		userService: userService,
		wsManager:   wsManager,
		clientMap:   clientMap,
	}
}

func (e *UserEventsExchange) Listen() {

	e.wsManager.AddObserver(events.User.AssociateClient, func(m websockets.IncomingWebSocketMessage, c *websockets.Client) {

		data := events.AssociateClientData{}
		if err := m.Payload.Assert(&data); err != nil {
			fmt.Println("Cannot type assert associate client data")
			return
		}

		claims, err := e.userService.VerifyUserToken(data.UserSecret)
		if err != nil {
			fmt.Println("VerifyUserToken error: ", err)
		} else {
			e.clientMap.Add(c.ID, claims.UserId)
			requestId, ok := m.Meta["rId"]
			if !ok {
				return
			}
			c.Emit(websockets.NewResponse(events.User.AssociateClient, requestId.(string), nil))
		}

	})

}
