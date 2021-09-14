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

package signalk

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/timmathews/argo/nmea2k"
)

type condition struct {
	Operation string `xml:"op"`
	Field     int    `xml:"field"`
	Value     string `xml:"value"`
}

type field struct {
	Key   string `xml:"type,attr"`
	Value int    `xml:",chardata"`
}

type fieldset struct {
	Type   string  `xml:"type,attr"`
	Fields []field `xml:"field"`
}

type parameterGroup struct {
	Pgn        uint32      `xml:"pgn"`
	Fieldset   fieldset    `xml:"fieldset"`
	Field      string      `xml:"field"`
	Multiplier multiplier  `xml:"multiplier"`
	Classifier classifier  `xml:"classifier"`
	Conditions []condition `xml:"condition"`
}

type sentence struct {
	Id    string `xml:"id"`
	Field string `xml:"field"`
}

type multiplier struct {
	Id    string `xml:"id"`
	Field string `xml:"field"`
}

type classifier struct {
	Id    string `xml:"id"`
	Field string `xml:"field"`
}

type mapping struct {
	Path            string           `xml:"path"`
	ParameterGroups []parameterGroup `xml:"parameter_group"`
	Sentences       []sentence       `xml:"sentence"`
}

type Mappings struct {
	XMLName  xml.Name  `xml:"mappings"`
	Mappings []mapping `xml:"mapping"`
}

type source struct {
	Pgn    uint32 `json:"pgn"`
	Device string `json:"device"`
	Src    uint8  `json:"src"`
}

type value struct {
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

type update struct {
	Source    source    `json:"source"`
	Timestamp time.Time `json:"timestamp"`
	Values    []value   `json:"values"`
}

type delta struct {
	Context string   `json:"context"`
	Updates []update `json:"updates"`
}

func ParseMappings(filename string) (Mappings, error) {
	output := Mappings{}

	data, err := ioutil.ReadFile(filename)

	if err == nil {
		err = xml.Unmarshal(data, &output)
	}

	return output, err
}

func (m *Mappings) Delta(msg *nmea2k.ParsedMessage) (delta, error) {
	context := "vessels.urn:mrn:signalk:uuid:c0d79334-4e25-4245-8892-54e8ccc8021d"

	src := source{
		Pgn:    msg.Header.Pgn,
		Device: "/dev/actisense",
		Src:    msg.Header.Source,
	}

	upd := update{
		Source:    src,
		Timestamp: time.Now(),
		Values:    *new([]value),
	}

	var usedFields = make(map[int]bool, len(msg.Data))

	for fieldId, field := range msg.Data {
		for _, mapping := range m.Mappings {
			for _, parameterGroup := range mapping.ParameterGroups {
				if parameterGroup.Pgn == msg.Header.Pgn {
					path := toDotNotation(mapping.Path)
					mlt := int64(-1)

					rpt := int64(nmea2k.PgnList[msg.Index].RepeatingFields)

					if parameterGroup.Classifier.Id != "" && parameterGroup.Classifier.Field != "" {
						cid, _ := strconv.ParseInt(parameterGroup.Classifier.Field, 10, 32)
						path = strings.Replace(path, fmt.Sprintf("{%v}", parameterGroup.Classifier.Id),
							fmt.Sprintf("%v", msg.Data[int(cid)]), 1)
					}

					if parameterGroup.Multiplier.Id != "" && parameterGroup.Multiplier.Field != "" {
						mlt, _ = strconv.ParseInt(parameterGroup.Multiplier.Field, 10, 32)

						sub := (((int64(fieldId) - mlt) / rpt) * rpt) + mlt

						path = strings.Replace(path, fmt.Sprintf("{%v}", parameterGroup.Multiplier.Id),
							fmt.Sprintf("%v", msg.Data[int(sub)]), 1)
					}

					fld, err := strconv.ParseInt(parameterGroup.Field, 10, 64)
					fid := int64(fieldId)

					if err == nil && ((fld > mlt && mlt > -1 && fld == fid%rpt) || fld == fid) && !usedFields[fieldId] {
						if conditionsMatch(parameterGroup.Conditions, msg.Data) {
							usedFields[fieldId] = true
							val := value{
								Path:  path,
								Value: field,
							}
							if val.Path != "" && val.Value != nil {
								upd.Values = append(upd.Values, val)
							}
						}
					} else if len(parameterGroup.Fieldset.Fields) > 0 {
						if parameterGroup.Fieldset.hasAll(msg.Data, usedFields) {
							s, u := parameterGroup.Fieldset.parse(msg.Data)
							val := value{
								Path:  toDotNotation(mapping.Path),
								Value: s,
							}
							if val.Path != "" && val.Value != nil {
								upd.Values = append(upd.Values, val)
							}
							usedFields = merge(usedFields, u)
						}
					}
				}
			}
		}
	}

	if len(upd.Values) > 0 {
		updates := []update{upd}
		delta := delta{
			Context: context,
			Updates: updates,
		}

		return delta, nil

	} else {
		return delta{}, fmt.Errorf("unknown PGN %v from %v on %v", upd.Source.Pgn, upd.Source.Src, upd.Source.Device)
	}
}

// Pack searches the mapping database for a matching path, then generates the PGN for that. This may
//func (m *Mappings) Pack(msg *update) (nmea2k.ParsedMessage, error) {
//
//}

// merge takes two maps and returns a new map containing the values of both
// maps. The right map takes precedence over the left map if both contain the
// same keys.
func merge(left, right map[int]bool) map[int]bool {
	var ret = make(map[int]bool, len(left)+len(right))

	for k, v := range left {
		ret[k] = v
	}

	for k, v := range right {
		ret[k] = v
	}

	return ret
}

// parse parses an intermediate format field map and returns a value and a map
// of used fields given a Fieldset. Currently parse only supports datetime type
// fieldsets.
func (fieldset fieldset) parse(dataFields nmea2k.DataMap) (string, map[int]bool) {
	f := fieldset.toMap()
	if fieldset.Type == "datetime" {
		dt, ok1 := dataFields[f["date"]].(time.Time)
		tm, ok2 := dataFields[f["time"]].(time.Time)

		if ok1 && ok2 {
			var usedFields = make(map[int]bool, 1)

			usedFields[f["date"]] = true
			usedFields[f["time"]] = true

			ts := tm.AddDate(dt.Year()-1970, int(dt.Month())-1, dt.Day()-1)
			return ts.Format(time.RFC3339), usedFields
		}
	}

	return "", nil
}

// toMap converts a Fieldset into a map.
func (fieldset fieldset) toMap() (out map[string]int) {
	out = make(map[string]int, 2)
	for _, field := range fieldset.Fields {
		out[field.Key] = field.Value
	}

	return
}

// hasAll checks if the incoming message contains all of the required fields to
// satisfy the conditions of a fieldset. This is generally only useful for date
// and time fields which need to be reassembled into a timestamp.
func (fieldset fieldset) hasAll(dataFields nmea2k.DataMap, usedFields map[int]bool) bool {
	for _, field := range fieldset.Fields {
		if usedFields[field.Value] || dataFields[field.Value] == "" || dataFields[field.Value] == nil {
			return false
		}
	}

	return true
}

// conditionsMatch iterates through an array of conditions and checks if the
// value of the referenced intermediate format field equates to the expected
// value using the operation provided.
func conditionsMatch(conditions []condition, fields nmea2k.DataMap) bool {
	for _, x := range conditions {
		ls := fmt.Sprintf("%v", fields[x.Field])

		switch x.Operation {
		case "eq":
			if ls != x.Value {
				return false
			}
		case "lt":
			if ls >= x.Value {
				return false
			}
		case "gt":
			if ls <= x.Value {
				return false
			}
		case "le":
			if ls > x.Value {
				return false
			}
		case "ge":
			if ls < x.Value {
				return false
			}
		default:
			return true
		}
	}

	return true
}

func toDotNotation(in string) string {
	return strings.Replace(in, "/", ".", -1)[2:]
}
