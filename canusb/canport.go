/*
 * Copyright (C) 2016 Tim Mathews <tim@signalk.org>
 *
 * This file is part of Argo.
 *
 * Argo is free software: you can redistribute it and/or modify it under the
 * terms of the GNU General Public License as published by the Free Software
 * Foundation, either version 3 of the License, or (at your option) any later
 * version.
 *
 * Argo is distributed in the hope that it will be useful, but WITHOUT ANY
 * WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
 * FOR A PARTICULAR PURPOSE. See the GNU General Public License for more
 * details.
 *
 * You should have received a copy of the GNU General Public License along with
 * this program. If not, see <http://www.gnu.org/licenses/>.
 */

package canusb

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/timmathews/argo/can"
)

type CanPort struct {
	p      io.ReadWriteCloser
	a      uint8
	IsOpen bool
}

var group byte = 0

// Send a 60928 ISO Address Claim parameter group
func (p *CanPort) AddressClaim(preferredAddress uint8) uint8 {
	buf := make([]byte, 8)
	unique := uint32(0x1fffff)
	manufacturer := uint32(100)
	lower_instance := uint32(0)
	upper_instance := uint32(0)
	function := uint32(25)
	class := uint32(25)
	instance := uint32(0)
	industry_code := uint32(4)
	arb_addr := uint32(1)

	var i0, i1 uint32

	i0 = unique
	i0 += manufacturer << 21

	i1 = lower_instance
	i1 += upper_instance << 3
	i1 += function << 8
	i1 += class << 17
	i1 += instance << 24
	i1 += industry_code << 28
	i1 += arb_addr << 31

	binary.LittleEndian.PutUint32(buf[0:4], i0)
	binary.LittleEndian.PutUint32(buf[4:8], i1)

	addr_claim := can.RawMessage{
		Timestamp:   time.Now(),
		Priority:    2,
		Source:      221,
		Destination: 225,
		Pgn:         60928,
		Length:      8,
		Data:        buf,
	}

	p.a = preferredAddress
	p.Send(&addr_claim)

	return preferredAddress
}

// OpenChannel opens the CAN bus port of the CANUSB adapter for communication.
// This must be called after opening the serial port, but before beginning
// communication with the CAN bus network. No harm will come from calling this
// function multiple times. CloseChannel is its counterpart.
func OpenChannel(port io.ReadWriteCloser, address uint8) (p *CanPort, err error) {
	defer func() {
		if err != nil && p != nil {
			p.CloseChannel()
		}
	}()

	// Set baudrate
	_, err = port.Write([]byte("S5\r")) // S5 = 250k
	if err != nil {
		return nil, err
	}

	// Open CANbus
	_, err = port.Write([]byte("O\r"))
	if err != nil {
		return nil, err
	}

	p = &CanPort{
		p:      port,
		IsOpen: false,
	}

	p.a = p.AddressClaim(address)
	p.IsOpen = true

	return p, nil
}

// CloseChannel closes the CAN bus port of the CANUSB adapter.
// This should be called before ending the communication session and must be
// called before closing the serial port. No harm will come from calling this
// function multiple times. OpenChannel is its counterpart.
func (p *CanPort) CloseChannel() error {
	_, err := p.Write([]byte("C\r"))

	p.p.Close()

	return err
}

func (p *CanPort) Read() (frame *can.RawMessage, err error) {
	rxbuf := []byte{0}
	msg := []byte{0}
	sof := false

	if p.IsOpen {
		for {
			_, err := p.p.Read(rxbuf)
			if err != nil {
				return nil, err
			}
			for _, b := range rxbuf {
				if b == 't' || b == 'T' || b == 'r' || b == 'R' {
					msg = nil
					msg = append(msg, b)
					sof = true
				} else if b == '\r' && sof {
					rec, err := p.frameReceived(msg)
					if err == nil {
						return &can.RawMessage{
							Timestamp:   time.Now(),
							Priority:    rec.Priority,
							Pgn:         rec.Pgn,
							Source:      rec.Source,
							Destination: rec.Destination,
							Length:      rec.Length,
							Data:        rec.Data,
						}, nil
					}
				} else if sof {
					msg = append(msg, b)
				}
			}
		}
	} else {
		return nil, errors.New("canusb.Read: CAN port is closed")
	}
}

func (p *CanPort) Address() uint8 {
	return p.a
}

func (p *CanPort) Write(b []byte) (int, error) {
	data := "T"
	pri := b[0]
	pgn := b[1:4]
	dst := b[4]
	len := b[5]
	pld := b[6 : 6+len]

	if len > 8 {
		return 0, errors.New("does not support long writes")
	}

	data += fmt.Sprintf("%.2X", pri<<2+pgn[0]&0x1)
	if pgn[1] < 240 {
		pgn[2] = dst
	}

	for _, byt := range pgn[1:] {
		data += fmt.Sprintf("%.2X", byt)
	}

	data += fmt.Sprintf("%.2X", p.a) // Source

	data += fmt.Sprintf("%.1X", len)

	for _, byt := range pld {
		data += fmt.Sprintf("%.2X", byt)
	}

	data += "\r"

	return p.p.Write([]byte(data))
}

// Send writes a RawMessage to the CANbus. Send can handles single-frame and
// fast packet PGNs. Larger data sets must be sent using
func (p *CanPort) Send(frame *can.RawMessage) (int, error) {
	buf := make([]byte, 14)

	buf[0] = frame.Priority
	buf[1] = byte((frame.Pgn & 0xf0000) >> 16)
	buf[2] = byte((frame.Pgn & 0xff00) >> 8)
	buf[3] = byte(frame.Pgn)
	buf[4] = frame.Destination

	// Up to eight bytes can be sent in a single frame
	// Up to 223 bytes can be sent as a fast packet
	// Over 223 bytes can be sent using other CANbus methods ... not implemented here

	dataLen := len(frame.Data)

	if dataLen <= 8 {
		buf[5] = frame.Length
		n := copy(buf[6:], frame.Data)

		if n < int(frame.Length) {
			return n, errors.New("not enough space for data")
		}

		return p.Write(buf)
	}

	if dataLen > 8 && dataLen <= 223 {
		chunksize := 6
		tmp := make([]byte, 8)
		seq := 0
		grp := group
		total := 0
		group++

		for i := 0; i < dataLen; i += chunksize {
			tmp[0] = byte(seq&0x1f) + byte(grp&0x7)<<5
			seq++

			if i == 1 {
				chunksize++
			}

			end := i + chunksize
			this := chunksize

			if end > dataLen {
				this = chunksize - (end - dataLen)
				end = dataLen
			}

			if i == 0 {
				tmp[1] = byte(dataLen)
				copy(tmp[2:], frame.Data[i:end])
				copy(buf[6:], tmp[0:this+2])
				buf[5] = byte(this + 2)

			} else {
				copy(tmp[1:], frame.Data[i:end])
				copy(buf[6:], tmp[0:this+1])
				buf[5] = byte(this + 1)
			}

			n, err := p.Write(buf)
			if err != nil {
				return n, err
			}

			total += n
		}

		return total, nil
	}

	return 0, errors.New("cannot send data larger than fast packets")
}

func (p *CanPort) frameReceived(msg []byte) (*CanFrame, error) {
	frame, err := ParseFrame(msg)
	if err != nil {
		return nil, err
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

	if isFastPacket(frame.Pgn) {
		frame.seq = frame.Data[0] & 0x1F
		frame.grp = (frame.Data[0] & 0x70) >> 5

		// PGN, source and group ID make a unique identifier for the frame group
		uid := uint32(uint32(frame.grp)<<28) +
			uint32(frame.Pgn<<8) +
			uint32(frame.Source)

		if frame.seq == 0 { // First in the series
			// Delete any existing scraps, should probably warn
			delete(partial_messages, uid)
			frame.Length = frame.Data[1]
			frame.Data = frame.Data[2:]

			if len(frame.Data) >= int(frame.Length) {
				return frame, nil
			} else {
				partial_messages[uid] = *frame
				return nil, errors.New("partial PGN")
			}
		} else {
			partial, ok := partial_messages[uid]
			if ok && partial.seq+1 == frame.seq {
				partial.Data = append(partial.Data, frame.Data[1:]...)
				partial.seq = frame.seq
				if len(partial.Data) >= int(partial.Length) {
					delete(partial_messages, uid)
					return &partial, nil
				} else {
					partial_messages[uid] = partial
					return nil, errors.New("partial PGN")
				}
			} // If we have a frame out of sequence, should probably warn
		}
	}

	return frame, nil
}
