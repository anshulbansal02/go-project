package roomaggregates

import "fmt"

func getUserToRoomRelationKey() string {
	return "relation:user->room"
}

func getRoomToUsersRelationKey(roomId string) string {
	return fmt.Sprintf("relation:room->users:%v", roomId)
}
