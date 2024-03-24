package main

import (
	"anshulbansal02/scribbly/controllers/middlewares"
	web "anshulbansal02/scribbly/controllers/web/http"
	exchange "anshulbansal02/scribbly/controllers/web/websocket"
	"anshulbansal02/scribbly/internal/chat"
	"anshulbansal02/scribbly/internal/room"
	"anshulbansal02/scribbly/internal/user"
	"anshulbansal02/scribbly/pkg/repository"
	"anshulbansal02/scribbly/pkg/websockets"

	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

func main() {

	// Load config
	// Setup loggers

	// Application Components Initialization
	repository := repository.New(&repository.Config{
		ServerAddress: "localhost:6379",
		Password:      "",
		DB:            0,
	})
	wsManager := websockets.NewWebSocketManager()
	rootRouter := chi.NewRouter()

	// Services Initialization
	userService := user.SetupConcreteService(*repository, user.Config{Secret: []byte("abce"), SigningMethod: jwt.SigningMethodHS256})
	roomService := room.SetupConcreteService(*repository)
	roomService.SetDependencies(room.DependingServices{
		UserService: userService,
	})
	chatService := chat.SetupConcreteService(*repository)

	rootRouter.Use(middlewares.Authenticate(userService))
	rootRouter.Use(middlewares.Cors)

	// Http Controllers Initialization
	rootRouter.Mount("/", web.SetupHealthcheckHttpControllers().Routes())
	rootRouter.Mount("/users", web.SetupUserHttpControllers(*userService).Routes())
	rootRouter.Mount("/rooms", web.SetupRoomHttpControllers(*roomService).Routes())

	rootRouter.HandleFunc("/client", wsManager.HandleWSConnection)

	// Events Exchange Setup
	clientMap := exchange.NewClientMap()
	exchange.NewRoomEventsExchange(roomService, wsManager, clientMap).Listen()
	exchange.NewUserEventsExchange(userService, wsManager, clientMap).Listen()
	exchange.NewChatEventsExchange(chatService, wsManager, clientMap).Listen()

	http.ListenAndServe("localhost:5000", rootRouter)

}
