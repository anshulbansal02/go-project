package events

import (
	"anshulbansal02/scribbly/pkg/websockets"
)

var User = struct {
	AssociateClient websockets.Event
}{
	AssociateClient: "associate_client",
}

type AssociateClientData struct {
	UserSecret string
}
