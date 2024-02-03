package room

import (
	"anshulbansal02/scribbly/pkg/utils"
	"fmt"
)

var generateRoomId = utils.NewRandomStringGenerator(utils.CHARSET_URL_SAFE, 8)
var generateRoomCode = utils.NewRandomStringGenerator(utils.CHARSET_ALPHA_NUM, 6)

func GetNamespaceKey(roomId string) string {
	return fmt.Sprintf("entity:room:%v", roomId)
}
