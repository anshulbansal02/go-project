package user

import "fmt"

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func getNamespaceKey(userId string) string {
	return fmt.Sprintf("entity:user:%v", userId)
}
