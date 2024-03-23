package chat

import (
	"anshulbansal02/scribbly/pkg/utils"

	"fmt"
)

var generateChatId = utils.NewRandomStringGenerator(utils.CHARSET_URL_SAFE, 12)

func getNamespaceKey(chatId string) string {
	return fmt.Sprintf("ent:chat:%v", chatId)
}
