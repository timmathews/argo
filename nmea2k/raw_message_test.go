package nmea2k

import (
	"math"
	"testing"
	//        "fmt"
	"time"
)

func TestDecodeLatLonWithValid64BitVal(t *testing.T) {
	vLat := 74.2343603
	data := make([]byte, 8)

	v := math.Float64bits(vLat)

	for i := 0; i < 8; i++ {
		data[i] = byte(v >> uint(i*8))
	}

	if x := decodeLatLon(RES_LATITUDE, data); x != vLat {
		t.Errorf("decodeLatLon(RES_LATITUDE, %v) = %v, expected %v", data, x, vLat)
	}
}

func TestDecodeLatLonWithValid32BitVal(t *testing.T) {
	var vLat float32
	vLat = 74.2343
	data := make([]byte, 4)

	v := math.Float32bits(vLat)

	for i := 0; i < 4; i++ {
		data[i] = byte(v >> uint(i*8))
	}

	if x := decodeLatLon(RES_LATITUDE, data); x != float64(vLat) {
		t.Errorf("decodeLatLon(RES_LATITUDE, %v) = %v, expected %v", data, x, float64(vLat))
	}
}

func TestDecodeDateWithValidTime(t *testing.T) {
	data := []byte{100, 0}

	tm := time.Date(1970, time.April, 10, 19, 0, 0, 0, time.Local)

	if x := decodeDate(data); x != tm {
		t.Errorf("decodeDate(%v) = %v, expected %v", data, x, tm)
	}
}

func TestDecodeTimeWithValidTime(t *testing.T) {
	data := []byte{0xFF, 0x97, 0x7F, 0x33}

	tm := time.Date(1970, time.January, 1, 23, 59, 59, 99990000, time.Local)

	if x := decodeTime(data); x != tm {
		t.Errorf("decodeTime(%v) = %v, expected %v", data, x, tm)
	}
}

func TestDecodeTemperatureWithValidTemp(t *testing.T) {
	data := []byte{0x91, 0xC3}

	temp := uint16(data[0]) | uint16(data[1])<<8

	temperature := float32(temp) / 100.0

	if x, _ := decodeTemperature(data); x != temperature {
		t.Errorf("decodeTemperature(%v) = %v, expected %v", data, x, temperature)
	}
}

func TestDecodePressureWithValidPressure(t *testing.T) {
	data := []byte{0x91, 0xC3}

	temp := uint16(data[0]) | uint16(data[1])<<8

	pressure := float32(temp) / 1000.0

	if x, _ := decodePressure(data); x != pressure {
		t.Errorf("decodePressure(%v) = %v, expected %v", data, x, pressure)
	}
}

func TestExtractNumber(t *testing.T) {

	data := []byte{0x06, 0xF0}

	startBit := uint32(3)
	bits := uint32(8)

	res := 4096

	if x := extractNumber(data, startBit, bits, 0); x != uint64(res) {
		t.Errorf("extractNumber(%v, %v, %v) = %v, expected %v", data, startBit, bits, x, res)
	}
}
