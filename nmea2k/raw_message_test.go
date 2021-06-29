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
	"encoding/hex"
	"testing"
	"time"

	"github.com/timmathews/argo/can"
)

type numTest struct {
	startByte uint32
	endByte   uint32
	startBit  uint32
	numBits   uint32
	expected  uint64
	data      []byte
}

func TestFieldOffsets(t *testing.T) {
	pgn := Pgn{"Test", "Test", 123456, true, 64, 0, []Field{
		{"Field 0", 8, 1, false, "", "", "", 0},  // Aligned, 1 byte
		{"Field 1", 12, 1, false, "", "", "", 0}, // Aligned, >1 byte
		{"Field 2", 4, 1, false, "", "", "", 0},  // Unaligned, <1 byte, no cross
		{"Field 3", 21, 1, false, "", "", "", 0}, // Aligned, >2 bytes
		{"Field 4", 5, 1, false, "", "", "", 0},  // Unaligned, <1 byte, cross
		{"Field 5", 14, 1, false, "", "", "", 0}, // Unaligned, >1 byte
		{"Field 6", 3, 1, false, "", "", "", 0},  // Aligned, <1 byte
		{"Field 7", 9, 1, false, "", "", "", 0},  // Unaligned, >1 byte
		{"Field 8", 20, 1, false, "", "", "", 0}, // Unaligned, >2 bytes
	}}

	expected := []numTest{
		{0, 1, 0, 8, 0, nil},
		{1, 3, 0, 12, 0, nil},
		{2, 3, 4, 4, 0, nil},
		{3, 6, 0, 21, 0, nil},
		{5, 6, 5, 5, 0, nil},
		{6, 8, 2, 14, 0, nil},
		{8, 9, 0, 3, 0, nil},
		{8, 10, 3, 9, 0, nil},
		{9, 12, 4, 20, 0, nil},
	}

	for i := 0; i < 9; i++ {
		r_start_byte, r_bytes, r_start_bit, r_bits := pgn.FieldOffsets(int32(i))
		e_values := expected[i]

		if !(r_start_byte == e_values.startByte && r_bytes == e_values.endByte &&
			r_start_bit == e_values.startBit && r_bits == e_values.numBits) {
			t.Errorf("FieldOffsets(%v) = %v, %v, %v, %v. Expected %v, %v, %v, %v",
				i, r_start_byte, r_bytes, r_start_bit, r_bits, e_values.startByte,
				e_values.endByte, e_values.startBit, e_values.numBits)
		}
	}
}

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
	// start byte, bytes, start bit, bits, expected, data
	data := []numTest{
		// Aligned to byte boundaries, whole byte sizes
		numTest{0, 1, 0, 8, 4, []byte{4}},
		numTest{0, 2, 0, 16, 4, []byte{4, 0}},
		numTest{0, 3, 0, 24, 4, []byte{4, 0, 0}},
		numTest{0, 4, 0, 32, 4, []byte{4, 0, 0, 0}},
		numTest{0, 5, 0, 40, 4, []byte{4, 0, 0, 0, 0}},
		numTest{0, 6, 0, 48, 4, []byte{4, 0, 0, 0, 0, 0}},
		numTest{0, 7, 0, 56, 4, []byte{4, 0, 0, 0, 0, 0, 0}},
		numTest{0, 8, 0, 64, 4, []byte{4, 0, 0, 0, 0, 0, 0, 0}},

		// Aligned to byte boundaries, fractional byte sizes
		numTest{0, 1, 0, 3, 1, []byte{0x79}},
		numTest{0, 2, 0, 11, 321, []byte{0x41, 0x79}},
		numTest{0, 3, 0, 22, 1655134, []byte{0x5e, 0x41, 0xD9}},
		numTest{0, 4, 0, 29, 155278871, []byte{0x17, 0x5e, 0x41, 0xC9}},
		numTest{0, 5, 0, 37, 39751391043, []byte{0x43, 0x17, 0x5e, 0x41, 0xC9}},

		// Real world values, a NAME
		numTest{0, 3, 0, 21, 110737, []byte{0x91, 0xb0, 0x21, 0x22, 0x00, 0x82, 0x32, 0xc0}},
		numTest{2, 4, 5, 11, 273, []byte{0x91, 0xb0, 0x21, 0x22, 0x00, 0x82, 0x32, 0xc0}},
		numTest{4, 5, 0, 3, 0, []byte{0x91, 0xb0, 0x21, 0x22, 0x00, 0x82, 0x32, 0xc0}},
		numTest{4, 5, 3, 5, 0, []byte{0x91, 0xb0, 0x21, 0x22, 0x00, 0x82, 0x32, 0xc0}},
		numTest{5, 6, 0, 8, 130, []byte{0x91, 0xb0, 0x21, 0x22, 0x00, 0x82, 0x32, 0xc0}},
		numTest{6, 7, 0, 1, 0, []byte{0x91, 0xb0, 0x21, 0x22, 0x00, 0x82, 0x32, 0xc0}},
		numTest{6, 7, 1, 7, 25, []byte{0x91, 0xb0, 0x21, 0x22, 0x00, 0x82, 0x32, 0xc0}},
		numTest{7, 8, 0, 4, 0, []byte{0x91, 0xb0, 0x21, 0x22, 0x00, 0x82, 0x32, 0xc0}},
		numTest{7, 8, 4, 3, 4, []byte{0x91, 0xb0, 0x21, 0x22, 0x00, 0x82, 0x32, 0xc0}},
		numTest{7, 8, 7, 1, 1, []byte{0x91, 0xb0, 0x21, 0x22, 0x00, 0x82, 0x32, 0xc0}},
	}

	field := Field{"", 0, 1, false, nil, "", "", 0}

	for _, d := range data {
		msg := RawMessage{new(can.RawMessage)}
		msg.Data = d.data

		if x, err := msg.extractNumber(&field, d.startByte, d.endByte, d.startBit, d.numBits); x != d.expected {
			t.Errorf("{%v}.extractNumber(1, %v, %v, %v, %v) = %v, expected %v",
				hex.EncodeToString(d.data), d.startByte, d.endByte, d.startBit, d.numBits, x, d.expected)
			t.Error(err)
		}
	}
}
