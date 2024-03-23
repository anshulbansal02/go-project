package user

import (
	"anshulbansal02/scribbly/pkg/utils"

	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

var generateUserId = utils.NewRandomStringGenerator(utils.CHARSET_URL_SAFE, 12)

func getNamespaceKey(userId string) string {
	return fmt.Sprintf("ent:user:%v", userId)
}

type UserClaims struct {
	UserId string `json:"userId"`
	jwt.RegisteredClaims
}
