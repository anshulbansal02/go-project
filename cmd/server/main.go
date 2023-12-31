package server

import (
	web "anshulbansal02/scribbly/controllers/web/http"
	"anshulbansal02/scribbly/internal/repository"
	"anshulbansal02/scribbly/internal/room"
	"anshulbansal02/scribbly/internal/user"
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

	// Services Initialization
	userService := user.NewService(user.NewRepository(*repository))
	roomService := room.NewService(room.NewRepository(*repository))
	// chatService := chat.NewService(chat.NewRepository(*repository))

	rootRouter := chi.NewRouter()
	// Http Controllers Initialization
	rootRouter.Mount("/users", web.SetupUserHttpControllers(*userService).Routes())
	rootRouter.Mount("/rooms", web.SetupRoomHttpControllers(*roomService).Routes())
	// rootRouter.Mount("/chat", web.SetupChatHttpControllers(*chatService).Routes)

	http.ListenAndServe(":6000", rootRouter)

}
