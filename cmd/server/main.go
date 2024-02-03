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

	rootRouter.Use(Cors)

	// Services Initialization
	userService := user.SetupConcreteService(*repository, user.Config{Secret: []byte("abce"), SigningMethod: jwt.SigningMethodHS256})
	roomService := room.SetupConcreteService(*repository)
	roomService.SetDependencies(room.DependingServices{
		UserService: userService,
	})

	// Http Controllers Initialization
	rootRouter.Mount("/", web.SetupHealthcheckHttpControllers().Routes())
	rootRouter.Mount("/users", web.SetupUserHttpControllers(*userService).Routes())
	rootRouter.Mount("/rooms", web.SetupRoomHttpControllers(*roomService).Routes())

	rootRouter.HandleFunc("/client", wsManager.HandleWSConnection)

	// Events Exchange Setup
	exchange.NewRoomEventsExchange(roomService, wsManager).Listen()
	exchange.NewUserEventsExchange(userService, wsManager).Listen()

	http.ListenAndServe(":5000", rootRouter)

}

// [TODO] Move to dedicated middleware directory
func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			return
		}
		next.ServeHTTP(w, r)
	})
}
