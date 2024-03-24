package room

import (
	"anshulbansal02/scribbly/pkg/utils"

	"fmt"
)

var generateRoomId = utils.NewRandomStringGenerator(utils.CHARSET_URL_SAFE, 12)
var generateRoomCode = utils.NewRandomStringGenerator(utils.CHARSET_ALPHA_UPPER+utils.CHARSET_NUM, 6)

func GetNamespaceKey(roomId string) string {
	return fmt.Sprintf("ent:room:%v", roomId)
}
