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
	"encoding/json"
	"fmt"
	"github.com/timmathews/argo/can"
	"gopkg.in/vmihailenco/msgpack.v2"
	"strconv"
	"time"
)

var layout = "2006-01-02T15:04:05.000"

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
	Timestamp   string                 `json:"timestamp"`
	Priority    uint8                  `json:"prio"`
	Source      uint8                  `json:"src"`
	Destination uint8                  `json:"dst"`
	Pgn         uint32                 `json:"pgn"`
	Fields      map[string]interface{} `json:"fields"`
}

func FromCanBoat(data string) (*ParsedMessage, error) {
	var cbm CanBoatMessage
	json.Unmarshal(([]byte)(data), &cbm)

	hdr := RawMessage{new(can.RawMessage)}

	hdr.Timestamp, _ = time.Parse("2006-01-02T15:04:05.999", cbm.Timestamp)
	hdr.Priority = cbm.Priority
	hdr.Pgn = cbm.Pgn
	hdr.Source = cbm.Source
	hdr.Destination = cbm.Destination

	f, fpgn := PgnList.First(cbm.Pgn)
	l, _ := PgnList.Last(cbm.Pgn)

	if f != l {
		return nil, fmt.Errorf("(%v): f (%v) != l (%v)\n", cbm.Pgn, f, l)
	}

	dd := make(map[int]interface{})

	for k, v := range cbm.Fields {
		for kk, vv := range fpgn.FieldList {
			if vv.Name == k {
				dd[kk] = v
			}
		}
	}

	var p = ParsedMessage{
		hdr,
		f,
		dd,
	}

	return &p, nil
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
