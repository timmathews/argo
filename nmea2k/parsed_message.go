package nmea2k

import (
	"encoding/json"
	"fmt"
	"gopkg.in/vmihailenco/msgpack.v2"
	"strconv"
	"time"
)

type DataMap map[int]interface{}

func (inVal DataMap) MarshalJSON() ([]byte, error) {
	outVal := make(map[string]interface{})

	for k, v := range inVal {
		outVal[strconv.Itoa(k)] = v
	}

	return json.Marshal(outVal)
}

type ParsedMessage struct {
	Header RawMessage
	Index  int
	Data   DataMap
}

type CanBoatMessage struct {
	Timestamp   time.Time `json:"timestamp"`
	Priority    uint8     `json:"prio"`
	Source      uint8     `json:"src"`
	Destination uint8     `json:"dst"`
	Pgn         uint32    `json:"pgn"`
	Fields      DataMap   `json:"fields"`
}

func FromCanBoat(data string) *ParsedMessage {
	var cbm CanBoatMessage
	json.Unmarshal(([]byte)(data), &cbm)
	fmt.Println("CBM", cbm)

	hdr := new(RawMessage)
	hdr.Timestamp = cbm.Timestamp
	hdr.Priority = cbm.Priority
	hdr.Pgn = cbm.Pgn
	hdr.Source = cbm.Source
	hdr.Destination = cbm.Destination

	i, _ := PgnList.First(cbm.Pgn)

	var p = ParsedMessage{
		*hdr,
		i,
		cbm.Fields,
	}

	return &p
}

func (msg *ParsedMessage) Print(verbose bool) string {
	// Timestamp Priority Source Destination Pgn PgnName: FieldName = FieldValue; ...

	pp := PgnList[msg.Index]
	name := pp.Description
	pgnFields := pp.FieldList

	s := fmt.Sprintf("%s %v %v %v %v %s:", msg.Header.Timestamp.Format(layout), msg.Header.Priority, msg.Header.Source, msg.Header.Destination, msg.Header.Pgn, name)

	for i, j := 0, 0; i < len(msg.Data); i, j = i+1, j+1 {
		f := msg.Data[i]
		if f != nil {
			if _, ok := f.(float32); ok {
				s += fmt.Sprintf(" %v.%s = %f;", i, pgnFields[j].Name, f)
			} else if _, ok := f.(float64); ok {
				s += fmt.Sprintf(" %v.%s = %f;", i, pgnFields[j].Name, f)
			} else {
				s += fmt.Sprintf(" %v.%s = %v;", i, pgnFields[j].Name, f)
			}
		} else if verbose {
			s += fmt.Sprintf(" %v.%s = nil;", i, pgnFields[j].Name)
		}

		if pp.RepeatingFields > 0 && j == len(pgnFields)-1 {
			j -= int(pp.RepeatingFields)
		}
	}

	s = s[:len(s)-1]

	return s
}

// Pack a PGN into a MsgPack formatted byte array
func (msg *ParsedMessage) MsgPack() []byte {
	b, err := msgpack.Marshal(&msg)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return b
}

// Pack a PGN into a JSON formatted byte array
func (msg *ParsedMessage) JSON() []byte {
	b, err := json.Marshal(&msg)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return b
}
