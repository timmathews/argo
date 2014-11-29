package can

import (
	"time"
)

type RawMessage struct {
	Timestamp   time.Time // Timestamp of receipt of CAN Message
	Priority    uint8     // Message priority, 0 is highest priority [id >> 26]
	Pgn         uint32    // Parameter group number/name id [(id & 0x3FFFFFF) >> 8]
	Source      uint8     // Sender ID [id & 0xFF]
	Destination uint8     /* Destination of message. If the PF field (bits 24:17 of the
	* ID) are >= 0xF0, than dst is 255 (broadcast). Otherwise
	* use the FS field (bits 16:9 of the ID) as the destination
	* address */
	Length uint8 // number of bytes which make up the frame
	Data   []byte
}

type Closer interface {
	CloseChannel() error
}

type Reader interface {
	Read() (*RawMessage, error)
}

type Writer interface {
	Write([]byte) (int, error)
}

type ReadWriter interface {
	Reader
	Writer
}
