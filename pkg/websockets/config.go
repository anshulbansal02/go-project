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
	messageTypeSectionBits     = 4  // Message type represented by 4 bits
	reservedBits               = 4  // Unused bits
	eventNameLengthSectionBits = 8  // Event Name Length represented by 8 bits or 1 byte (event name can be of max 256 bytes/characters)
	metaDataLengthSectionBits  = 16 // Meta Data Length represented by 16 bites or 2 bytes (meta data can be of max 65,536 bytes/characters)
	metaBytes                  = (messageTypeSectionBits +
		reservedBits +
		eventNameLengthSectionBits +
		metaDataLengthSectionBits) / 8 // Total size of the meta bytes in a message or minimum size of the message
	maxEventNameLength = 1 << eventNameLengthSectionBits
)
