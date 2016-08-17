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

package main

import (
	"encoding/xml"
	"fmt"
	"github.com/timmathews/argo/nmea2k"
	"strconv"
	"strings"
	"time"
)

const timeformat = "2006-01-02T15:04:05"

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
	Pgn       uint32    `json:"pgn"`
	Device    string    `json:"device"`
	Timestamp time.Time `json:"timestamp"`
	Src       uint8     `json:"src"`
}

type value struct {
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

type update struct {
	Source source  `json:"source"`
	Values []value `json:"values"`
}

type delta struct {
	Context string   `json:"context"`
	Updates []update `json:"updates"`
}

func GenerateMappings(mappingData []byte) (Mappings, error) {
	output := Mappings{}

	err := xml.Unmarshal(mappingData, &output)

	return output, err
}

func (m *Mappings) Delta(msg *nmea2k.ParsedMessage) (update, error) {
	src := source{
		Pgn:       msg.Header.Pgn,
		Device:    "/dev/actisense",
		Timestamp: time.Now(),
		Src:       msg.Header.Source,
	}

	update := update{
		Source: src,
		Values: *new([]value),
	}

	var usedFields = make(map[int]bool, len(msg.Data))

	for fieldId, field := range msg.Data {
		for _, mapping := range m.Mappings {
			for _, parameterGroup := range mapping.ParameterGroups {
				if parameterGroup.Pgn == msg.Header.Pgn {
					path := mapping.Path
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
								update.Values = append(update.Values, val)
							}
						}
					} else if len(parameterGroup.Fieldset.Fields) > 0 {
						if parameterGroup.Fieldset.hasAll(msg.Data, usedFields) {
							s, u := parameterGroup.Fieldset.parse(msg.Data)
							val := value{
								Path:  mapping.Path,
								Value: s,
							}
							if val.Path != "" && val.Value != nil {
								update.Values = append(update.Values, val)
							}
							usedFields = merge(usedFields, u)
						}
					}
				}
			}
		}
	}

	if len(update.Values) > 0 {
		return update, nil
	} else {
		return update, fmt.Errorf("unknown PGN %v from %v on %v", update.Source.Pgn, update.Source.Src, update.Source.Device)
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
		dt, ok := dataFields[f["date"]].(time.Time)
		tm, ok := dataFields[f["time"]].(time.Time)

		if ok {
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
