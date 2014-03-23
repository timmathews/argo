package actisense

import (
	"argo/nmea2k"
	"fmt"
	"github.com/schleibinger/sio"
	"log"
	"time"
)

const (
	// ASCII characters which mark packet start and stop
	STX = 0x02
	ETX = 0x03
	DLE = 0x10
	ESC = 0x1B

	// N2K commands
	N2K_MSG_RECEIVED = 0x93
	N2K_MSG_SEND     = 0x94

	// NGT commands
	NGT_MSG_RECEIVED = 0xA0
	NGT_MSG_SEND     = 0xA1
)

type MsgState int

const (
	MSG_START MsgState = iota
	MSG_ESCAPE
	MSG_MESSAGE
)

/*
 * Wrap the PGN or NGT message and send to NGT
 *
 * The message envelope has the following structure:
 *
 * <DLE><STX><COMMAND><LEN><CMD DATA><CRC><DLE><ETX>
 *
 * <COMMAND> is a one byte to either send or receive a specific
 * N2K or NGT message
 *
 * <LEN> is the length of the unescaped <CMD DATA>
 *
 * <CMD DATA> is the actual command being sent, either an NGT message or an
 * NMEA2000 PGN. Any DLE characters (0x10) are escaped with another DLE
 * character, so <DLE> becomes <DLE><DLE>.
 *
 * <CRC> is such that the sum of all unescaped data bytes plus the command byte
 * plus the length plus the checksum add up to zero, modulo 256.
 */
func WriteMessage(port *sio.Port, command byte, payload []byte) {

	bst := []byte{DLE, STX}

	bst = append(bst, command, byte(len(payload)))

	crc := command

	for _, b := range payload {
		if b == DLE {
			bst = append(bst, DLE)
		}
		bst = append(bst, b)
		crc += b
	}

	crc += byte(len(payload))

	crc = byte(256 - int(crc))

	bst = append(bst, crc, DLE, ETX)

	n, err := port.Write(bst)

	if err != nil {
		log.Fatalln("write: %s", err)
	}

	if n != len(bst) {
		log.Fatalf("short write: %d %d", n, len(bst))
	}

	log.Printf("Wrote command %v len %d\n", bst, len(bst))
}

func ReadNGT1(port *sio.Port, data chan byte, result chan nmea2k.ParsedMessage) {
	var buf []byte
	state := MSG_START
	for b := range data {
		if state == MSG_ESCAPE {
			if b == ETX { // End of message
				messageReceived(buf, result)
				buf = nil
				state = MSG_START
			} else if b == STX { // Start of message
				state = MSG_MESSAGE
			} else if b == DLE { // Escaped DLE char
				buf = append(buf, b)
				state = MSG_MESSAGE
			} else { // Unexpected character after DLE
				buf = nil
				state = MSG_START
			}
		} else if state == MSG_MESSAGE {
			if b == DLE { // Escape char
				state = MSG_ESCAPE
			} else {
				buf = append(buf, b)
			}
		} else {
			if b == DLE { // Escape char
				state = MSG_ESCAPE
			}
		}
	}
}

func messageReceived(msg []byte, res chan nmea2k.ParsedMessage) {

	if len(msg) < 3 {
		fmt.Printf("Ignore short command len = %v\n", len(msg))
		return
	}

	var checksum byte
	for _, c := range msg {
		checksum += c
	}

	if checksum != 0 {
		fmt.Printf("Ignoring message with invalid checksum")
		return
	}

	command := msg[0]

	//fmt.Printf("Message command = %v, %d\n", msg, len(msg))

	if command == N2K_MSG_RECEIVED {
		n2kMessageReceived(msg[1:], res)
	} else if command == NGT_MSG_RECEIVED {
		ngtMessageReceived(msg[1:], res)
	} else {
		fmt.Printf("Unknown message type (%02X) received", command)
	}
}

func n2kMessageReceived(msg []byte, res chan nmea2k.ParsedMessage) {

	// Packet length from NGT1
	if msg[0] < 11 {
		log.Println("Ignore short msg", len(msg))
		return
	}

	raw := new(nmea2k.RawMessage)
	raw.Timestamp = time.Now()
	raw.Priority = msg[1]
	raw.Pgn = uint32(msg[2]) | uint32(msg[3])<<8 | uint32(msg[4])<<16
	raw.Destination = msg[5]
	raw.Source = msg[6]
	// Skip the timestamp (bytes 7-10)
	lth := msg[11]

	if lth > 223 {
		log.Println("Ignore long msg", lth)
		return
	}

	raw.Length = lth
	lth += 12
	raw.Data = msg[12:lth]

	parsed := raw.ParsePacket()

	res <- *parsed
}

func ngtMessageReceived(msg []byte, res chan nmea2k.ParsedMessage) {

	pLen := msg[0]

	if pLen < 12 {
		log.Println("Ignore short msg", len(msg))
		return
	}

	raw := new(nmea2k.RawMessage)
	raw.Timestamp = time.Now()
	raw.Priority = 0
	raw.Pgn = 0x40000 + uint32(msg[1])
	raw.Destination = 0
	raw.Source = 0
	raw.Length = pLen - 1
	pLen++
	raw.Data = msg[2:pLen]

	parsed := raw.ParsePacket()

	res <- *parsed
}
