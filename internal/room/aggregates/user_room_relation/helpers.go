package roomaggregates

import "fmt"

func getUserToRoomRelationKey() string {
	return "rel:user->room"
}

func getRoomToUsersRelationKey(roomId string) string {
	return fmt.Sprintf("rel:room->users:%v", roomId)
}
