package canusb

import (
	"errors"
	"fmt"
	"github.com/schleibinger/sio"
	"github.com/timmathews/argo/nmea2k"
	"log"
	"strconv"
	"time"
)

type msgType int

const (
	CAN_STD msgType = iota
	CAN_EXT
	CAN_STD_RTR
	CAN_EXT_RTR
)

type CanFrame struct {
	msgType msgType // Standard or extended or request message
	id      uint32  // Full ID of frame, may be removed in future releases
	pgn     uint32  // Parameter group number/name id [(id & 0x3FFFFFF) >> 8]
	pri     uint8   // Message priority, 0 is highest priority [id >> 26]
	src     uint8   // Sender ID [id & 0xFF]
	dst     uint8   /* Destination of message. If the PF field (bits 24:17 of the
								   * ID) are >= 0xF0, than dst is 255 (broadcast). Otherwise
	                 * use the FS field (bits 16:9 of the ID) as the destination
	                 * address */
	length uint8 // number of bytes which make up the frame
	grp    uint8 // group for fast packet
	seq    uint8 // sequence of the frame in a fast packet
	data   []byte
}

func (frm *CanFrame) String() string {
	str := fmt.Sprintf("%v: %d %d %d %d %d %d: ", frm.msgType, frm.pri, frm.src,
		frm.dst, frm.pgn, frm.length, frm.seq)

	for _, b := range frm.data {
		str += fmt.Sprintf("[%.2x]", b)
	}

	str += "\n"

	return str
}

// Storage for list of fast packet PGNs
var fast_packets = make([]uint32, 0)

// Storage for incomplete, received fast packets
var partial_messages = make(map[uint32]CanFrame)

var portIsOpen = false

// OpenChannel opens the CAN bus port of the CANUSB adapter for communication.
// This must be called after opening the serial port, but before beginning
// communication with the CAN bus network. No harm will come from calling this
// function multiple times. CloseChannel is its counterpart.
func OpenChannel(port *sio.Port) {
	var s string

	// Set baudrate
	s = fmt.Sprintf("S%d\r", 5) // 5 = 250k
	_, err := port.Write([]byte(s))
	if err != nil {
		log.Fatalln("Failed to set speed")
	}

	// Open CANbus
	s = fmt.Sprintf("O\r")
	_, err = port.Write([]byte(s))
	if err != nil {
		log.Fatalln("Failed to open CANbus")
	}

	portIsOpen = true
}

// CloseChannel closes the CAN bus port of the CANUSB adapter.
// This should be called before ending the communication session and must be
// called before closing the serial port. No harm will come from calling this
// function multiple times. OpenChannel is its counterpart.
func CloseChannel(port *sio.Port) {
	var s string

	fmt.Sprintf(s, "C\r")
	_, err := port.Write([]byte(s))
	if err != nil {
		log.Fatalln("Failed to close CANbus")
	}

	portIsOpen = false
}

// Adds fast packet PGNs to a unique, sorted array
func AddFastPacket(pgn uint32) {
	for i, v := range fast_packets {
		if pgn == v {
			return
		} else if pgn < v {
			fast_packets = append(fast_packets, 0)
			copy(fast_packets[i+1:], fast_packets[i:])
			fast_packets[i] = pgn
			return
		}
	}
	fast_packets = append(fast_packets, pgn)
}

func isFastPacket(pgn uint32) bool {
	i, j := 0, len(fast_packets)
	for i < j {
		h := i + (j-i)/2
		if !(fast_packets[h] >= pgn) {
			i = h + 1
		} else {
			j = h
		}
	}

	return i < len(fast_packets) && fast_packets[i] == pgn
}

func ParseFrame(p []byte) (*CanFrame, error) {

	frame := new(CanFrame)
	var n, v uint64
	var err error
	var offset int

	switch p[0] {
	case 't':
		{
			frame.msgType = CAN_STD
			offset = 4
		}
	case 'T':
		{
			frame.msgType = CAN_EXT
			offset = 9
		}
	case 'r':
		{
			frame.msgType = CAN_STD_RTR
			offset = 4
		}
	case 'R':
		{
			frame.msgType = CAN_EXT_RTR
			offset = 9
		}
	default:
		return nil, errors.New("canusb.ParseFrame: Invalid prefix")
	}

	n, err = strconv.ParseUint(string(p[1:offset]), 16, offset*4)
	if err == nil {
		frame.id = uint32(n)
		frame.pri = uint8(frame.id >> 26)
		frame.pgn = (frame.id & 0x3FFFFFF) >> 8
		frame.src = uint8(frame.id & 0xFF)
		pf := (frame.id & 0xFFFFFF) >> 16
		if pf >= 240 {
			frame.dst = 255
		} else {
			frame.dst = uint8(pf)
		}
	} else {
		return nil, errors.New(fmt.Sprintf("canusb.ParseFrame: Unable to parse message ID: %s", err))
	}

	n, err = strconv.ParseUint(string(p[offset]), 16, 8)
	if err == nil {
		if n <= 8 {
			frame.length = uint8(n)
		} else {
			return nil, errors.New("canusb.ParseFrame: Expected length <= 8")
		}
	} else {
		return nil, errors.New(fmt.Sprintf("canusb.ParseFrame: Unable to parse message length: %s", err))
	}

	offset++

	data_len := len(p[offset:]) - 4

	if data_len%2 != 0 || data_len/2 != int(n) {
		return nil, errors.New(fmt.Sprintf("canusb.ParseFrame: Expected %d bytes, got %d", n*2, data_len))
	} else {
		for i := offset; i < data_len+offset; i += 2 {
			v, err = strconv.ParseUint(string(p[i:i+2]), 16, 8)
			if err == nil {
				frame.data = append(frame.data, byte(v))
			} else {
				return nil, errors.New(fmt.Sprintf("canusb.ParseFrame: Unable to parse data: %s", err))
			}
		}
	}

	return frame, nil
}

func Write(port *sio.Port, payload *nmea2k.RawMessage) {
	if payload.Destination < 255 {
		payload.Pgn += uint32(payload.Destination)
	}
	addr := uint32(payload.Priority<<26) + uint32(payload.Pgn<<8) + uint32(payload.Source)
	data := fmt.Sprintf("T%.8X%.1X", addr, payload.Length)
	for _, b := range payload.Data {
		data += fmt.Sprintf("%.2X", b)
	}
	data += "/r"
	port.Write([]byte(data))
}

func ReadPort(data chan byte, result chan nmea2k.ParsedMessage) {
	msg := []byte{0}
	sof := false

	if portIsOpen {
		for b := range data {
			if b == 't' || b == 'T' || b == 'r' || b == 'R' {
				msg = nil
				msg = append(msg, b)
				sof = true
			} else if b == '\r' && sof == true {
				frameReceived(msg, result)
				msg = nil
				sof = false
			} else if sof == true {
				msg = append(msg, b)
			}
		}
	}
}

func frameReceived(msg []byte, result chan nmea2k.ParsedMessage) {

	frame, err := ParseFrame(msg)
	if err != nil {
		log.Fatalln(err)
	}

	// data[0] bits 7-5: group ID ... i.e. all of these belong together, unless
	//                   time between packets exceeds an unknown number of ms
	// data[0] bits 4-0: sequence ... the number of this frame in the sequence
	//                   of fast packet frames. since we do not know if packets
	//                   are allowed out of order, assume it is not allowed
	// data[1]: if sequence is 0 this is the total number of bytes in the fast
	//          packet set, otherwise it is part of the data
	//
	// As a result of the conditions above, fast packets can be up to 223 bytes.
	// 5 bits for sequence means up to 32 total frames in a fast packet. A frame
	// can have at most 8 bytes of data, but in fast packet mode the first byte
	// is always group ID and sequence. Also the first frame of a fast packet can
	// only have 6 bytes because the second byte is the byte count for the packet
	//
	// 223 = 31 * 7 + 6
	//
	// Should we bail if we see a byte count > 223?

	if isFastPacket(frame.pgn) {

		frame.seq = frame.data[0] & 0x1F
		frame.grp = (frame.data[0] & 0x70) >> 5

		// PGN, source and group ID make a unique identifier for the frame group
		uid := uint32(frame.grp<<28) + uint32(frame.pgn<<8) + uint32(frame.src)

		if frame.seq == 0 { // First in the series
			delete(partial_messages, uid) // Delete any existing scraps, should probably warn
			frame.length = frame.data[1]
			frame.data = frame.data[2:]

			if len(frame.data) >= int(frame.length) {
				messageReceived(frame, result)
			} else {
				partial_messages[uid] = *frame
			}
		} else {
			p, ok := partial_messages[uid]
			if ok && p.seq+1 == frame.seq {
				p.data = append(p.data, frame.data[1:]...)
				p.seq = frame.seq
				if len(p.data) >= int(p.length) {
					messageReceived(&p, result)
					delete(partial_messages, uid)
				} else {
					partial_messages[uid] = p
				}
			} // If we have a frame out of sequence, should probably warn
		}
	} else {
		messageReceived(frame, result)
	}
}

func messageReceived(frame *CanFrame, res chan nmea2k.ParsedMessage) {
	raw := new(nmea2k.RawMessage)

	raw.Timestamp = time.Now()
	raw.Priority = frame.pri
	raw.Pgn = frame.pgn
	raw.Destination = frame.dst
	raw.Source = frame.src

	raw.Data = frame.data

	parsed := raw.ParsePacket()

	res <- *parsed
}
