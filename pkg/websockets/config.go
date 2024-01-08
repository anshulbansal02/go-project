package websockets

import "time"

var (
	newline = []byte{'\n'}
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1000000 // 1 MB
)

const (
	messageTypeSectionBits         = 4
	reservedBits                   = 4
	eventNameLengthSectionBits     = 8
	metaBits                       = messageTypeSectionBits + reservedBits + eventNameLengthSectionBits
	maxEventNameLength         int = 1 << eventNameLengthSectionBits
)
