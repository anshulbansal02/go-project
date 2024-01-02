package user

import (
	"anshulbansal02/scribbly/pkg/utils"
	"fmt"
)

var generateUserId = utils.NewRandomStringGenerator(nil, 12)

func getNamespaceKey(userId string) string {
	return fmt.Sprintf("entity:user:%v", userId)
}
