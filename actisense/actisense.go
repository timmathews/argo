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

package actisense

import (
	"errors"
	"fmt"
	"github.com/schleibinger/sio"
	"github.com/timmathews/argo/can"
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

/* The following startup command reverse engineered from Actisense NMEAreader.
 * It instructs the NGT1 to clear its PGN message TX list, thus it starts
 * sending all PGNs.
 */
var NGT_STARTUP_SEQ = []byte{0x11, 0x02, 0x00}

type MsgState int

const (
	MSG_START MsgState = iota
	MSG_ESCAPE
	MSG_MESSAGE
)

type ActisensePort struct {
	p      *sio.Port
	IsOpen bool
}

func OpenChannel(port *sio.Port) (p *ActisensePort, err error) {
	p = &ActisensePort{
		p:      port,
		IsOpen: true,
	}

	_, err = p.write(NGT_MSG_SEND, NGT_STARTUP_SEQ)

	if err != nil {
		return nil, err
	}

	return p, nil
}

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
func (p *ActisensePort) Write(payload []byte) (int, error) {
	return p.write(N2K_MSG_SEND, payload)
}

func (p *ActisensePort) Read() (*can.RawMessage, error) {
	var buf []byte
	rxbuf := []byte{0}
	state := MSG_START
	var msg *can.RawMessage

	for {
		_, err := p.p.Read(rxbuf)
		if err != nil {
			return nil, err
		}

		for _, b := range rxbuf {
			if state == MSG_ESCAPE {
				if b == ETX { // End of message
					msg, err = messageReceived(buf)
					buf = nil
					state = MSG_START
					if err == nil {
						return msg, nil
					}
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
}

func (p *ActisensePort) write(command byte, payload []byte) (int, error) {
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

	return p.p.Write(bst)
}

func messageReceived(msg []byte) (*can.RawMessage, error) {

	if len(msg) < 3 {
		return nil, errors.New(fmt.Sprintf("Ignore short command len = %v\n", len(msg)))
	}

	var checksum byte
	for _, c := range msg {
		checksum += c
	}

	if checksum != 0 {
		return nil, errors.New("Ignoring message with invalid checksum")
	}

	command := msg[0]

	if command == N2K_MSG_RECEIVED {
		return n2kMessageReceived(msg[1:])
	} else if command == NGT_MSG_RECEIVED {
		return ngtMessageReceived(msg[1:])
	} else {
		return nil, errors.New(fmt.Sprintf("Unknown message type (%02X) received", command))
	}
}

func n2kMessageReceived(msg []byte) (*can.RawMessage, error) {

	// Packet length from NGT1
	if msg[0] < 11 {
		return nil, errors.New(fmt.Sprintf("Ignore short msg", len(msg)))
	}

	raw := new(can.RawMessage)
	raw.Timestamp = time.Now()
	raw.Priority = msg[1]
	raw.Pgn = uint32(msg[2]) | uint32(msg[3])<<8 | uint32(msg[4])<<16
	raw.Destination = msg[5]
	raw.Source = msg[6]
	// Skip the timestamp (bytes 7-10)
	lth := msg[11]

	if lth > 223 {
		return nil, errors.New(fmt.Sprintf("Ignore long msg", lth))
	}

	raw.Length = lth
	lth += 12
	raw.Data = msg[12:lth]

	return raw, nil
}

func ngtMessageReceived(msg []byte) (*can.RawMessage, error) {

	pLen := msg[0]

	if pLen < 12 {
		return nil, errors.New(fmt.Sprintf("Ignore short msg", len(msg)))
	}

	raw := new(can.RawMessage)
	raw.Timestamp = time.Now()
	raw.Priority = 0
	raw.Pgn = 0x40000 + uint32(msg[1])
	raw.Destination = 0
	raw.Source = 0
	raw.Length = pLen - 1
	pLen++
	raw.Data = msg[2:pLen]

	return raw, nil
}
