package nmea2k

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

const layout = "2006-01-02-15:04:05.999"

type DecodeError struct {
	Data  []byte
	Where string
}

func (e *DecodeError) Error() string {
	return fmt.Sprintf("%v is not valid data for %s", e.Data, e.Where)
}

type RawMessage struct {
	Timestamp   time.Time
	Priority    uint8
	Source      uint8
	Destination uint8
	Pgn         uint32
	Length      uint8
	Data        []byte
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

func (msg *RawMessage) ParsePacket() (pgnParsed *ParsedMessage) {

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

	for idx := 0; idx < len(fields); idx++ {

		field := fields[idx]
		res := field.Resolution

		bits := field.Size
		bytes := (bits + 7) / 8
		bytes = min(bytes, uint32(data_len))
		bits = min(bytes*8, bits)

		bytes += start_byte

		var data interface{}
		var err error

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
				data, err = msg.extractNumber(res, start_byte, bytes, start_bit, bits)
			case RES_LOOKUP:
				data, err = msg.extractLookupField(&field, start_byte, bytes, start_bit, bits)
			case RES_MANUFACTURER:
				data, err = msg.extractManufacturer(&field, start_byte, bytes, start_bit, bits)
			case RES_PRESSURE:
				data, err = msg.extractPressure(start_byte, bytes)
			case RES_STRINGLZ:
				data, err = msg.extractString(start_byte)
			default:
				data = msg.Data[start_byte:bytes]
			}
		} else if field.Resolution > 0.0 {
			data, err = msg.extractNumber(res, start_byte, bytes, start_bit, bits)
		}

		if err == nil {
			pgnParsed.Data[idx] = data

			if !oneSolution {
				for i <= j {
					if res == RES_MANUFACTURER {
						data, _ = msg.extractNumber(res, start_byte, bytes, start_bit, bits)
					}
					if v, ok := field.Units.(string); ok && v[0] == '=' {
						value, _ := strconv.ParseUint(v[1:], 10, 64)
						if value == data.(uint64) {
							// fmt.Println(idx, pgnParsed.Data[idx], value)
							// We have a match, continue parsing
							break
						} else {
							i++
							pgnDefinition = PgnList[i]
							fields = pgnDefinition.FieldList
							field = fields[idx]
						}
					} else {
						break
					}
				}
				pgnParsed.Index = i
			}
		} else {
			pgnParsed.Data[idx] = nil
		}

		start_byte = start_byte + ((bits + start_bit) / 8)
		start_bit += bits
		start_bit %= 8

	}

	return
}

func (msg *RawMessage) extractLatLon(start, end uint32) (v float64, e error) {
	var value uint64

	data := msg.Data[start:end]
	bytes := len(data)

	for i := 0; i < bytes; i++ {
		value |= uint64(data[i]) << uint(8*i)
	}

	if value == 0xFFFF || value == 0xFFFFFFFF {
		v = math.NaN()
		e = &DecodeError{data, "Data not present"}
	}

	if bytes == 4 {
		v = float64(math.Float32frombits(uint32(value)))
	} else if bytes == 8 {
		v = math.Float64frombits(value)
	} else {
		v = math.NaN()
		e = &DecodeError{data, "Invalid float"}
	}

	return

}

// Date values are 16 bits and represent number of days since the Unix Epoch
func (msg *RawMessage) extractDate(start, end uint32) (t time.Time, e error) {
	var d uint32

	data := msg.Data[start:end]

	if len(data) != 2 {
		e = &DecodeError{data, "Field size mismatch"}
	}

	d = uint32(data[0]) | uint32(data[1])<<8

	if d == 0xFFFF {
		e = &DecodeError{data, "Data not present"}
	}

	t = time.Unix(int64(d*86400), 0)

	return

}

func (msg *RawMessage) extractTime(start, end uint32) (t time.Time, e error) {
	var d uint32

	data := msg.Data[start:end]

	if len(data) != 4 {
		e = &DecodeError{data, "Field size mismatch"}
	}

	for i := 0; i < 4; i++ {
		d |= uint32(data[i]) << uint(8*i)
	}

	if d == 0xFFFFFFFF {
		e = &DecodeError{data, "Data not present"}
	}

	seconds := d / 10000
	units := d % 10000
	minutes := seconds / 60
	seconds = seconds % 60
	hours := minutes / 60
	minutes = minutes % 60

	t = time.Date(1970, time.January, 1, int(hours), int(minutes), int(seconds), int(units*10000), time.Local)

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

func (msg *RawMessage) extractString(start uint32) (s string, e error) {

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

	if int(end) >= len(msg.Data) {
		e = &DecodeError{nil, "Data not present"}
		return
	}

	start++

	data := msg.Data[start:end]

	if len(data) == 0 {
		e = &DecodeError{data, "Data not present"}
	} else {
		s = string(data)
	}

	return

}

func (msg *RawMessage) extractNumber(res float64, start, end, startBit, bits uint32) (value interface{}, e error) {

	data := msg.Data[start:end]
	bytes := len(data)

	if bytes > 8 {
		e = &DecodeError{data, "Field size mismatch"}
	}

	var num uint64
	var nbits, shift uint32
	var t byte

	for i := 0; i < bytes; i++ {
		shift += nbits
		// Calculate number of bits used in the byte
		if i == 0 {
			nbits = min(8-startBit, bits)
		} else if (i == bytes-1) && (bits%8 != 0) {
			nbits = min(8, (bits-startBit)%8)
		} else {
			nbits = 8
		}

		mask := byte(0xFF)
		mask = mask >> byte(8-nbits)

		t = data[i] >> startBit

		num |= uint64(t&mask) << shift

	}

	var maxValue uint64
	if bits > 8 {
		maxValue = 1<<(bits-1) - 1
	} else {
		maxValue = 1<<bits - 1
	}

	if maxValue == 0 {
		maxValue = 0x7FFFFFFFFFFFFFFF
	}

	if num >= maxValue {
		e = &DecodeError{data, "Field not present"}
		return
	}

	if res != 1 && res != RES_LOOKUP && res != RES_MANUFACTURER && res != RES_INTEGER {
		value = float64(num) * float64(res)
	} else {
		value = num
	}

	return

}

func (msg *RawMessage) extractLookupField(f *Field, start, end, startBit, bits uint32) (ret interface{}, e error) {

	res := f.Resolution

	n, err := msg.extractNumber(res, start, end, startBit, bits)

	if err != nil {
		e = err
		return
	}

	if _, ok := f.Units.(map[int]string); ok {
		v := f.Units.(map[int]string)[int(n.(uint64))]
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

func (msg *RawMessage) extractManufacturer(f *Field, start, end, startBit, bits uint32) (ret interface{}, e error) {

	res := f.Resolution

	n, err := msg.extractNumber(res, start, end, startBit, bits)

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

func (msg *RawMessage) Print(verbose bool) (s string) {
	// Timestamp Priority Source Destination Pgn PgnName: Data
	s = fmt.Sprintf("%s,%v,%v,%v,%v,%v,% x", msg.Timestamp.Format(layout), msg.Priority, msg.Pgn, msg.Source, msg.Destination, msg.Length, msg.Data)
	return
}

func (msg *RawMessage) GetPgnDefinition(pgn uint32) *Pgn {
	_, p := PgnList.First(msg.Pgn)
	return &p
}
