package nmea2k

import (
	"encoding/json"
	"fmt"
	msgpack "github.com/vmihailenco/msgpack"
	"strconv"
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

func (msg *ParsedMessage) Print(verbose bool) string {
	// Timestamp Priority Source Destination Pgn PgnName: FieldName = FieldValue; ...

	pp := PgnList[msg.Index]
	name := pp.Description
	pgnFields := pp.FieldList

	s := fmt.Sprintf("%s %v %v %v %v %s:", msg.Header.Timestamp.Format(layout), msg.Header.Priority, msg.Header.Source, msg.Header.Destination, msg.Header.Pgn, name)

	for i := 0; i < len(msg.Data); i++ {
		f := msg.Data[i]
		if f != nil {
			if _, ok := f.(float32); ok {
				s += fmt.Sprintf(" %v.%s = %f;", i, pgnFields[i].Name, f)
			} else if _, ok := f.(float64); ok {
				s += fmt.Sprintf(" %v.%s = %f;", i, pgnFields[i].Name, f)
			} else {
				s += fmt.Sprintf(" %v.%s = %v;", i, pgnFields[i].Name, f)
			}
		} else if verbose {
			s += fmt.Sprintf(" %v.%s = nil;", i, pgnFields[i].Name)
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
