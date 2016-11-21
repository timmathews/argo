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

package can

import (
	"fmt"
	"time"
)

const layout = "2006-01-02-15:04:05.999"

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

func (msg *RawMessage) Print(verbose bool) (s string) {
	// Timestamp Priority Pgn Source Destination Length Data
	s = fmt.Sprintf("%s %v %v %v %v %v: % x", msg.Timestamp.Format(layout), msg.Priority, msg.Source, msg.Destination, msg.Pgn, msg.Length, msg.Data)
	return
}
