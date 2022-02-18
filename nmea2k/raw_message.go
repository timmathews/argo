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
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/timmathews/argo/can"
)

type DecodeError struct {
	Data  []byte
	Where string
}

type RawMessage struct {
	*can.RawMessage
}

func (e *DecodeError) Error() string {
	return fmt.Sprintf("%v is not valid data for %s", e.Data, e.Where)
}

func min(x, y uint32) uint32 {
	if x < y {
		return x
	}

	return y
}

func max(x, y uint32) uint32 {
	if x > y {
		return x
	}

	return y
}

func (p *Pgn) FieldOffsets(idx int32) (low_byte, high_byte, start_bit, bits uint32) {
	bits = p.FieldList[idx].Size
	bytes := (bits + 7) / 8
	bits = min(bytes*8, bits)

	var offset uint32
	for i, f := range p.FieldList {
		if i == int(idx) {
			break
		}
		offset += f.Size
	}

	low_byte = offset / 8
	high_byte = low_byte + bytes
	start_bit = offset % 8

	return
}

func ParsePacket(cmsg *can.RawMessage) (pgnParsed *ParsedMessage) {
	msg := &RawMessage{cmsg}

	i, pgnDefinition := PgnList.First(msg.Pgn)
	j, _ := PgnList.Last(msg.Pgn)

	oneSolution := false

	if i == j {
		oneSolution = true
	}

	pgnParsed = new(ParsedMessage)
	pgnParsed.Header = *msg
	pgnParsed.Index = i
	pgnParsed.Data = make(map[int]interface{})

	data_len := len(msg.Data)

	fields := pgnDefinition.FieldList

	var start_byte uint32
	var start_bit uint32

	for idx, odx := 0, 0; idx < len(fields); idx, odx = idx+1, odx+1 {

		field := fields[idx]
		res := field.Resolution

		bits := field.Size
		bytes := (bits + 7) / 8
		bytes = min(bytes, uint32(data_len))
		bits = min(bytes*8, bits)

		bytes += start_byte

		var data interface{}
		var err error

		// TODO: We should return an error
		if int(start_byte) > len(msg.Data) || int(bytes) > len(msg.Data) {
			return
		}

		// Special fields
		if field.Resolution < 0.0 {
			switch field.Resolution {
			case RES_LATITUDE:
				fallthrough
			case RES_LONGITUDE:
				data, err = msg.extractLatLon(start_byte, bytes)
			case RES_DATE:
				data, err = msg.extractDate(start_byte, bytes)
			case RES_TIME:
				data, err = msg.extractTime(start_byte, bytes)
			case RES_TEMPERATURE:
				data, err = msg.extractTemperature(start_byte, bytes)
			case RES_6BITASCII:
				data, err = msg.extract6BitASCII(start_byte, bytes)
			case RES_INTEGER:
				data, err = msg.extractNumber(&field, start_byte, bytes, start_bit, bits)
			case RES_LOOKUP:
				data, err = msg.extractLookupField(&field, start_byte, bytes, start_bit, bits)
			case RES_LOOKUP2:
				superFieldId := field.Offset
				superFieldVal := pgnParsed.Data[int(superFieldId)]
				// The superField may not have been parsed yet
				if superFieldVal == nil {
					start_byte, bytes, start_bit, bits := pgnDefinition.FieldOffsets(superFieldId)
					data, err = msg.extractNumber(&field, start_byte, bytes, start_bit, bits)
					if err != nil {
						break
					}
				} else {
					data = uint32(superFieldVal.(float64))
				}
				data, err = msg.extractLookupSubfield(&field, uint32(data.(float64)), start_byte, bytes, start_bit, bits)
			case RES_MANUFACTURER:
				data, err = msg.extractManufacturer(&field, start_byte, bytes, start_bit, bits)
			case RES_PRESSURE:
				data, err = msg.extractPressure(start_byte, bytes)
			case RES_STRINGLZ:
				data, err = msg.extractStringLZ(start_byte)
			case RES_STRING:
				data = string(msg.Data[start_byte:bytes])
			case RES_ASCII:
				data, err = msg.extractString(start_byte, bytes)
			default:
				data = msg.Data[start_byte:bytes]
			}
		} else if field.Resolution > 0.0 {
			data, err = msg.extractNumber(&field, start_byte, bytes, start_bit, bits)
		}

		if err == nil {
			pgnParsed.Data[odx] = data

			if !oneSolution {
				for i <= j {
					if res == RES_MANUFACTURER {
						data, _ = msg.extractNumber(&field, start_byte, bytes, start_bit, bits)
					}
					if v, ok := field.Units.(string); ok && v[0] == '=' {
						value, _ := strconv.ParseUint(v[1:], 10, 64)
						if value == data.(uint64) {
							// We have a match, continue parsing
							break
						} else {
							i++
							pgnListLen := len(PgnList)
							if i < pgnListLen {
								pgnDefinition = PgnList[i]
								fields = pgnDefinition.FieldList

								if idx >= len(fields) {
									break
								}

								field = fields[idx]
							} else {
								break
							}
						}
					} else {
						break
					}
				}
				pgnParsed.Index = i
			}
		} else {
			pgnParsed.Data[odx] = nil
		}

		start_byte = start_byte + ((bits + start_bit) / 8)
		start_bit += bits
		start_bit %= 8

		if idx == len(fields)-1 && pgnDefinition.RepeatingFields > 0 && start_byte < uint32(data_len) {
			idx -= int(pgnDefinition.RepeatingFields)
		}
	}

	return
}

func (msg *RawMessage) extractLatLon(start, end uint32) (v interface{}, e error) {
	data := msg.Data[start:end]
	bytes := len(data)

	if bytes == 4 {
		var value int32
		for i, b := range data {
			value |= int32(b) << uint(8*i)
		}
		if value > 0x7FFFFFFD {
			v = math.NaN()
			e = &DecodeError{data, "Data not present"}
		} else {
			v = float32(value) / 1e+7
		}
	} else if bytes == 8 {
		var value int64
		for i, b := range data {
			value |= int64(b) << uint(8*i)
		}
		if value > 0x7FFFFFFFFFFFFFFD {
			v = math.NaN()
			e = &DecodeError{data, "Data not present"}
		} else {
			v = float64(value) / 1e+16
		}
	} else {
		v = math.NaN()
		e = &DecodeError{data, "Invalid float"}
	}

	if e != nil {
		fmt.Println(e)
	}
	return
}

// Date values are 16 bits and represent number of days since the Unix Epoch
func (msg *RawMessage) extractDate(start, end uint32) (t time.Time, e error) {
	var d uint32

	data := msg.Data[start:end]

	if len(data) != 2 {
		e = &DecodeError{data, "Field size mismatch"}
	} else {
		d = uint32(data[0]) | uint32(data[1])<<8
		if d == 0xFFFF {
			e = &DecodeError{data, "Data not present"}
		} else {
			t = time.Unix(int64(d*86400), 0)
		}
	}

	return
}

func (msg *RawMessage) extractTime(start, end uint32) (t time.Time, e error) {
	var d uint32

	data := msg.Data[start:end]

	if len(data) != 4 {
		e = &DecodeError{data, "Field size mismatch"}
	} else {
		for i := 0; i < 4; i++ {
			d |= uint32(data[i]) << uint(8*i)
		}
		if d == 0xFFFFFFFF {
			e = &DecodeError{data, "Data not present"}
		} else {
			seconds := d / 10000
			units := d % 10000
			minutes := seconds / 60
			seconds = seconds % 60
			hours := minutes / 60
			minutes = minutes % 60

			t = time.Date(1970, time.January, 1, int(hours), int(minutes), int(seconds), int(units*10000), time.Local)
		}
	}

	return

}

// Takes 2 bytes and returns a 32 bit float representing the temperature in
// Kelvin.
func (msg *RawMessage) extractTemperature(start, end uint32) (t float32, e error) {

	data := msg.Data[start:end]

	if len(data) != 2 {
		e = &DecodeError{data, "Field size mismatch"}
		return
	}

	d := uint16(data[0]) | uint16(data[1])<<8

	if d >= 0xfffd {
		e = &DecodeError{data, "Data not present"}
		return
	}

	t = float32(d) / 100.0

	return

}

func (msg *RawMessage) extract6BitASCII(start, end uint32) (s string, e error) {

	data := msg.Data[start:end]

	e = &DecodeError{data, "Not implemented"}

	return

}

// Takes 2 bytes and returns a 32 bit float representing the pressure in bar
func (msg *RawMessage) extractPressure(start, end uint32) (p float32, e error) {

	data := msg.Data[start:end]

	if len(data) != 2 {
		e = &DecodeError{data, "Field size mismatch"}
		return
	}

	d := uint16(data[0]) | uint16(data[1])<<8

	if d >= 0xfffd {
		e = &DecodeError{data, "Data not present"}
		return
	}

	p = float32(d) / 1000.0

	return

}

func (msg *RawMessage) extractStringLZ(start uint32) (s string, e error) {
	if int(start) >= len(msg.Data) {
		e = &DecodeError{nil, "Data not present"}
		return
	}

	end := msg.Data[start]
	if end == 0 {
		e = &DecodeError{nil, "Data not present"}
		return
	}

	end += byte(start)
	start++

	if int(end) > len(msg.Data) {
		e = &DecodeError{nil, "Data not present"}
		return
	}

	data := msg.Data[start:end]
	if len(data) == 0 {
		e = &DecodeError{data, "Data not present"}
	} else {
		s = string(data)
	}

	return
}

func (msg *RawMessage) extractString(start, end uint32) (s string, e error) {

	if int(start) >= len(msg.Data) || int(end) >= len(msg.Data) {
		e = &DecodeError{nil, "Data not present"}
		return
	}

	if msg.Data[start] == 0 {
		e = &DecodeError{nil, "Data not present"}
		return
	}

	i := start
	for ; i < end; i++ {
		if msg.Data[int(i)] == 0 || msg.Data[int(i)] == 255 {
			break
		}
	}

	data := msg.Data[start:i]

	if len(data) == 0 {
		e = &DecodeError{data, "Data not present"}
	} else {
		s = string(data)
	}

	return

}

func (msg *RawMessage) extractNumber(field *Field, start, end, offset, width uint32) (value interface{}, e error) {
	data := make([]byte, 8)
	bytes := end - start
	res := field.Resolution
	var num uint64

	if bytes > 8 {
		e = &DecodeError{msg.Data[start:end], "Numeric field exceeds max width"}
		return
	}

	copy(data, msg.Data[start:end])
	mask := ^uint64(0) >> uint64(64-width)

	num = binary.LittleEndian.Uint64(data)
	num = num >> offset
	num = num & mask

	if (1<<width)&num != 0 {
		e = &DecodeError{data, "Field not present"}
		return
	}

	var maxValue uint64
	if width > 8 {
		maxValue = 1<<(width-1) - 1
	} else if width == 1 {
		maxValue = 1
	} else {
		maxValue = 1<<width - 1
	}

	if maxValue == 0 {
		maxValue = 0x7FFFFFFFFFFFFFFF
	}

	if num > maxValue {
		e = &DecodeError{data, "Field not present"}
		return
	}

	if res != 1 && res != RES_LOOKUP && res != RES_MANUFACTURER && res != RES_INTEGER {
		if field.Signed {
			if field.Size <= 8 {
				value = float64(int8(num)) * float64(res)
			} else if field.Size <= 16 {
				value = float64(int16(num)) * float64(res)
			} else if field.Size <= 32 {
				value = float64(int32(num)) * float64(res)
			} else {
				value = float64(int64(num)) * float64(res)
			}
		} else {
			value = float64(num) * float64(res)
		}
	} else {
		if field.Signed {
			if field.Size <= 8 {
				value = int8(num)
			} else if field.Size <= 16 {
				value = int16(num)
			} else if field.Size <= 32 {
				value = int32(num)
			} else {
				value = int64(num)
			}
		} else {
			value = num
		}
	}

	return
}

func (msg *RawMessage) extractLookupField(f *Field, start, end, startBit, bits uint32) (ret interface{}, e error) {
	n, err := msg.extractNumber(f, start, end, startBit, bits)

	if err != nil {
		e = err
		return
	}

	if _, ok := f.Units.(PgnLookup); ok {
		v := f.Units.(PgnLookup)[int(n.(uint64))]
		if v == "" {
			ret = n
		} else {
			ret = v
		}
	} else {
		ret = n
	}

	return

}

func (msg *RawMessage) extractLookupSubfield(f *Field, superId, start, end, startBit, bits uint32) (ret interface{}, e error) {
	n, err := msg.extractNumber(f, start, end, startBit, bits)
	if err != nil {
		e = err
		return
	}

	if u, ok := f.Units.(PgnSubLookup); ok {
		ret = u[int(superId)][int(n.(float64))]
	} else {
		ret = n
	}

	return
}

func (msg *RawMessage) extractManufacturer(f *Field, start, end, startBit, bits uint32) (ret interface{}, e error) {
	n, err := msg.extractNumber(f, start, end, startBit, bits)

	if err != nil {
		e = err
		return
	}

	v := lookupCompanyCode[int(n.(uint64))]
	if v == "" {
		ret = n
	} else {
		ret = v
	}

	return

}

func (msg *RawMessage) GetPgnDefinition(pgn uint32) *Pgn {
	_, p := PgnList.First(msg.Pgn)
	return &p
}
