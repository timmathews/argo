package canusb

import (
	"errors"
	"fmt"
	"github.com/timmathews/argo/can"
	"strconv"
)

type msgType int

const (
	CAN_STD msgType = iota
	CAN_EXT
	CAN_STD_RTR
	CAN_EXT_RTR
)

// Storage for incomplete, received fast packets
var partial_messages = make(map[uint32]CanFrame)

// Storage for list of fast packet PGNs
var fast_packets = make([]uint32, 0)

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

type CanFrame struct {
	can.RawMessage
	msgType msgType // Standard or extended or request message
	id      uint32  // Full ID of frame, may be removed in future releases
	grp     uint8   // group for fast packet
	seq     uint8   // sequence of the frame in a fast packet
}

func (frm *CanFrame) String() string {
	str := fmt.Sprintf("%v: %d %d %d %d %d %d: ", frm.msgType, frm.Priority,
		frm.Source, frm.Destination, frm.Pgn, frm.Length, frm.seq)

	for _, b := range frm.Data {
		str += fmt.Sprintf("[%.2x]", b)
	}

	str += "\n"

	return str
}

func ParseFrame(p []byte) (*CanFrame, error) {

	frame := new(CanFrame)
	var n, v uint64
	var err error
	var offset int

	if p == nil || len(p) == 0 {
		return nil, errors.New("Empty byte array")
	}

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
		frame.Priority = uint8(frame.id >> 26)
		frame.Source = uint8(frame.id & 0xFF)
		pf := (frame.id >> 16) & 0xFF
		ps := (frame.id >> 8) & 0xFF
		if pf >= 240 {
			frame.Destination = 255
			frame.Pgn = (frame.id >> 8) & 0x03FFFF
		} else {
			frame.Destination = uint8(ps)
			frame.Pgn = (frame.id >> 8) & 0x03FF00
		}
	} else {
		return nil, errors.New(fmt.Sprintf("canusb.ParseFrame: Unable to parse message ID: %s", err))
	}

	n, err = strconv.ParseUint(string(p[offset]), 16, 8)
	if err == nil {
		if n <= 8 {
			frame.Length = uint8(n)
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
				frame.Data = append(frame.Data, byte(v))
			} else {
				return nil, errors.New(fmt.Sprintf("canusb.ParseFrame: Unable to parse data: %s", err))
			}
		}
	}

	return frame, nil
}
