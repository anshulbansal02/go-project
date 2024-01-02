package main

import (
	web "anshulbansal02/scribbly/controllers/web/http"
	exchange "anshulbansal02/scribbly/controllers/web/websocket"
	"anshulbansal02/scribbly/internal/repository"
	"anshulbansal02/scribbly/internal/room"
	"anshulbansal02/scribbly/internal/user"
	"anshulbansal02/scribbly/pkg/websockets"
	"net/http"

	"github.com/go-chi/chi/v5"
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
	userService := user.SetupConcreteService(*repository)
	roomService := room.SetupConcreteService(*repository)
	roomService.SetDependencies(room.DependingServices{
		UserService: userService,
	})
	// chatService := chat.NewService(chat.NewRepository(*repository))

	// Http Controllers Initialization
	rootRouter.Mount("/", web.SetupHealthcheckHttpControllers().Routes())
	rootRouter.Mount("/users", web.SetupUserHttpControllers(*userService).Routes())
	rootRouter.Mount("/rooms", web.SetupRoomHttpControllers(*roomService).Routes())
	// rootRouter.Mount("/chat", web.SetupChatHttpControllers(*chatService).Routes)

	// Events Exchange Setup
	exchange.NewRoomEventsExchange(roomService, wsManager).Listen()

	http.ListenAndServe(":6000", rootRouter)

}
