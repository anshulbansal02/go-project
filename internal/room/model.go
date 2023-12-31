package room

import "fmt"

type Room struct {
	ID           string   `json:"id"`
	Code         string   `json:"code"`
	Type         string   `json:"type"`
	Participants []string `json:"participants_ids"`
	Admin        *string  `json:"admin_id"`
}

func GetNamespaceKey(roomId string) string {
	return fmt.Sprintf("entity:room:%v", roomId)
}
