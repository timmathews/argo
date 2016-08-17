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

package nmea2k

import (
	"github.com/timmathews/argo/can"
	"testing"
	"time"
)

func TestExtractLatLonWithValid64BitVal(t *testing.T) {
	vLat := -76.56960748136044
	data := make([]byte, 8)
	msg := RawMessage{new(can.RawMessage)}

	v := int64(vLat * 1e16)

	for i := 0; i < 8; i++ {
		data[i] = byte(v >> uint(i*8))
	}

	msg.Data = data

	if x, err := msg.extractLatLon(0, 8); x != vLat {
		t.Errorf("decodeLatLon(RES_LATITUDE, %v) = %v, expected %v", data, x, vLat)
		t.Error(err)
	}
}

func TestExtractLatLonWithValid32BitVal(t *testing.T) {
	var vLat float32
	vLat = -76.5696
	data := make([]byte, 4)
	msg := RawMessage{new(can.RawMessage)}

	v := int32(vLat * 1e7)

	for i := 0; i < 4; i++ {
		data[i] = byte(v >> uint(i*8))
	}

	msg.Data = data

	if x, err := msg.extractLatLon(0, 4); x != vLat {
		t.Errorf("decodeLatLon(RES_LATITUDE, %v) = %v, expected %v", data, x, vLat)
		t.Error(err)
	}
}

func TestExtractDateWithValidTime(t *testing.T) {
	msg := RawMessage{new(can.RawMessage)}
	msg.Data = []byte{100, 0}

	tm := time.Date(1970, time.April, 10, 19, 0, 0, 0, time.Local)

	if x, err := msg.extractDate(0, 2); x != tm {
		t.Errorf("decodeDate(%v) = %v, expected %v", msg.Data, x, tm)
		t.Error(err)
	}
}

func TestExtractTimeWithValidTime(t *testing.T) {
	msg := RawMessage{new(can.RawMessage)}
	msg.Data = []byte{0xFF, 0x97, 0x7F, 0x33}

	tm := time.Date(1970, time.January, 1, 23, 59, 59, 99990000, time.Local)

	if x, err := msg.extractTime(0, 4); x != tm {
		t.Errorf("decodeTime(%v) = %v, expected %v", msg.Data, x, tm)
		t.Error(err)
	}
}

func TestExtractTemperatureWithValidTemp(t *testing.T) {
	msg := RawMessage{new(can.RawMessage)}
	msg.Data = []byte{0x91, 0xC3}

	temp := uint16(msg.Data[0]) | uint16(msg.Data[1])<<8

	temperature := float32(temp) / 100.0

	if x, err := msg.extractTemperature(0, 2); x != temperature {
		t.Errorf("decodeTemperature(%v) = %v, expected %v", msg.Data, x, temperature)
		t.Error(err)
	}
}

func TestExtractPressureWithValidPressure(t *testing.T) {
	msg := RawMessage{new(can.RawMessage)}
	msg.Data = []byte{0x91, 0xC3}

	temp := uint16(msg.Data[0]) | uint16(msg.Data[1])<<8

	pressure := float32(temp) / 1000.0

	if x, err := msg.extractPressure(0, 2); x != pressure {
		t.Errorf("decodePressure(%v) = %v, expected %v", msg.Data, x, pressure)
		t.Error(err)
	}
}

func TestExtractNumber(t *testing.T) {
	msg := RawMessage{new(can.RawMessage)}
	msg.Data = []byte{0x06, 0xF0}

	startBit := uint32(3)
	bits := uint32(8)

	res := 4096

	if x, err := msg.extractNumber(1, 0, 2, startBit, bits); x != uint64(res) {
		t.Errorf("extractNumber(%v, %v, %v) = %v, expected %v", msg.Data, startBit, bits, x, res)
		t.Error(err)
	}
}
