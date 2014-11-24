package canusb

import (
	"errors"
	"fmt"
	"github.com/schleibinger/sio"
	"log"
)

type CanPort struct {
	p      *sio.Port
	a      uint8
	IsOpen bool
	rx     chan []byte
	tx     chan *CanFrame
}

// OpenChannel opens the CAN bus port of the CANUSB adapter for communication.
// This must be called after opening the serial port, but before beginning
// communication with the CAN bus network. No harm will come from calling this
// function multiple times. CloseChannel is its counterpart.
func OpenChannel(port *sio.Port, address uint8) (p *CanPort, err error) {
	var s string

	defer func() {
		if err != nil && p != nil {
			port.Close()
		}
	}()

	// Set baudrate
	s = fmt.Sprintf("S%d\r", 5) // 5 = 250k
	_, err = port.Write([]byte(s))
	if err != nil {
		return nil, err
	}

	// Open CANbus
	s = fmt.Sprintf("O\r")
	_, err = port.Write([]byte(s))
	if err != nil {
		return nil, err
	}

	// TODO: Address negotiation
	p = &CanPort{
		p:      port,
		a:      221,
		IsOpen: true,
		rx:     make(chan []byte),
		tx:     make(chan *CanFrame),
	}

	go p.run()

	return p, nil
}

// CloseChannel closes the CAN bus port of the CANUSB adapter.
// This should be called before ending the communication session and must be
// called before closing the serial port. No harm will come from calling this
// function multiple times. OpenChannel is its counterpart.
func (p *CanPort) CloseChannel() error {
	var s string

	fmt.Sprintf(s, "C\r")
	_, err := p.Write([]byte(s))

	close(p.tx)
	close(p.rx)

	return err
}

func (p *CanPort) Read() (frame *CanFrame, err error) {
	rxbuf := []byte{0}
	msg := []byte{0}
	sof := false

	if p.IsOpen {
		_, err := p.p.Read(rxbuf)
		if err != nil {
			return nil, err
		}
		for _, b := range rxbuf {
			if b == 't' || b == 'T' || b == 'r' || b == 'R' {
				msg = nil
				msg = append(msg, b)
				sof = true
			} else if b == '\r' && sof == true {
				p.rx <- msg
				msg = nil
				sof = false
			} else if sof == true {
				msg = append(msg, b)
			}
		}
	} else {
		return nil, errors.New("canusb.Read: CAN port is closed")
	}

	frame = <-p.tx

	return frame, nil
}

func (p *CanPort) Write(b []byte) (int, error) {
	data := "T"
	for _, byt := range b {
		data += fmt.Sprintf("%.2X", byt)
	}
	data += "/r"
	return p.p.Write([]byte(data))
}

func (p *CanPort) run() {
	for {
		select {
		case c := <-p.rx:
			p.frameReceived(c)
		}
	}
}

func (p *CanPort) frameReceived(msg []byte) {

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

	if isFastPacket(frame.Pgn) {

		frame.seq = frame.Data[0] & 0x1F
		frame.grp = (frame.Data[0] & 0x70) >> 5

		// PGN, source and group ID make a unique identifier for the frame group
		uid := uint32(frame.grp<<28) + uint32(frame.Pgn<<8) + uint32(frame.Src)

		if frame.seq == 0 { // First in the series
			delete(partial_messages, uid) // Delete any existing scraps, should probably warn
			frame.Length = frame.Data[1]
			frame.Data = frame.Data[2:]

			if len(frame.Data) >= int(frame.Length) {
				p.tx <- frame
			} else {
				partial_messages[uid] = *frame
			}
		} else {
			partial, ok := partial_messages[uid]
			if ok && partial.seq+1 == frame.seq {
				partial.Data = append(partial.Data, frame.Data[1:]...)
				partial.seq = frame.seq
				if len(partial.Data) >= int(partial.Length) {
					p.tx <- &partial
					delete(partial_messages, uid)
				} else {
					partial_messages[uid] = partial
				}
			} // If we have a frame out of sequence, should probably warn
		}
	} else {
		p.tx <- frame
	}
}
