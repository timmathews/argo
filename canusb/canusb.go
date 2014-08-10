package canusb

import (
	"argo/nmea2k"
	"github.com/schleibinger/sio"
	"time"
	"fmt"
	"log"
	"errors"
	"strconv"
)

type msgType int

const (
	CAN_STD msgType = iota
	CAN_EXT
	CAN_STD_RTR
	CAN_EXT_RTR
)

type CanFrame struct {
	msgType msgType
	id			uint32
	pri			uint8
	pgn			uint32
	src			uint8
	dst			uint8
	length  uint8
	data    []byte
}

func (frm *CanFrame) String() string {
	str := fmt.Sprintln("Type:", frm.msgType)
	str += fmt.Sprintf("ID:  %x\n", frm.id)
	str += fmt.Sprintf("Pri: %x\n", frm.pri)
	str += fmt.Sprintf("PGN: %x\n", frm.pgn)
	str += fmt.Sprintf("Src: %x\n", frm.src)
	str += fmt.Sprintf("Dst: %x\n", frm.dst)
	str += fmt.Sprintln("Len: ", frm.length)
	str += fmt.Sprint("Data: ")
	for _, b := range frm.data {
		str += fmt.Sprintf("[%.2x]", b)
	}
	str += fmt.Sprintln()
	return str
}

var portIsOpen = false

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

func CloseChannel(port *sio.Port) {
	var s string

	fmt.Sprintf(s, "C\r");
	_, err := port.Write([]byte(s))
	if err != nil {
		log.Fatalln("Failed to close CANbus")
	}

	portIsOpen = false
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
		frame.pri = uint8(frame.id>>26)
		frame.pgn = (frame.id & 0x1FFFFFF)>>8
		frame.src = uint8(frame.id & 0xFF)
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

func WriteMessage(port *sio.Port, command byte, payload []byte) {
}

func ReadPort(data chan byte, result chan nmea2k.ParsedMessage) {
	msg := []byte{0}
	sof := false

	if portIsOpen {
		for b := range data {
			if b == 't' ||
			b == 'T' ||
			b == 'r' ||
			b == 'R' {
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

	//frame, err := ParseFrame(msg)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	messageReceived(msg, result)
}

func messageReceived(msg []byte, res chan nmea2k.ParsedMessage) {

	frame, err := ParseFrame(msg)
	if err != nil {
		log.Fatalln(err)
	}

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

