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
	set "github.com/deckarep/golang-set"
	"github.com/timmathews/argo/can"
	"github.com/timmathews/argo/nmea2k"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

var mapdata Mappings

func MakeSet(s []value) set.Set {
	a := set.NewSet()
	for _, item := range s {
		a.Add(item)
	}

	return a
}

func TestMain(m *testing.M) {
	mapfile, err := ioutil.ReadFile("./map.xml")
	if err != nil {
		log.Fatal(err)
	}

	err = xml.Unmarshal(mapfile, &mapdata)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func TestFieldsetValidDate(t *testing.T) {
	ts := time.Now()
	in := nmea2k.ParsedMessage{
		nmea2k.RawMessage{&can.RawMessage{
			ts,
			3, 126992, 1, 255, 8, []byte{0x0, 0xF, 0xC2, 0x40, 0xD0, 0x89, 0x00, 0x00},
		}},
		68,
		nmea2k.DataMap{0: 0, 1: "GPS", 2: 0xF, 3: time.Unix(16578*86400, 0).UTC(), 4: time.Unix(43200, 0).UTC()},
	}

	expected := update{
		source{126992, "/dev/actisense", ts, 1},
		[]value{{"~/system/currentTime", "2015-05-23T12:00:00Z"}, {"~/system/currentTimeSource", "GPS"}},
	}

	got, err := mapdata.Delta(&in)

	x := MakeSet(got.Values)
	y := MakeSet(expected.Values)

	if !x.Equal(y) {
		t.Errorf("\nExpected: %+v\n     Got: %+v\n     Err: %v", expected, got, err)
	} else {
		fmt.Printf("%+v\n", got)
	}
}

func TestFieldsetMissingDate(t *testing.T) {
	ts := time.Now()
	in := nmea2k.ParsedMessage{
		nmea2k.RawMessage{&can.RawMessage{
			ts,
			3, 126992, 1, 255, 8, []byte{0x0, 0xF, 0xC2, 0x40, 0xD0, 0x89, 0x00, 0x00},
		}},
		68,
		nmea2k.DataMap{0: 0, 1: "GPS", 2: 0xF, 4: time.Unix(43200, 0).UTC()},
	}

	expected := update{
		source{126992, "/dev/actisense", ts, 1},
		[]value{{"~/system/currentTimeSource", "GPS"}},
	}

	got, err := mapdata.Delta(&in)

	x := MakeSet(got.Values)
	y := MakeSet(expected.Values)

	if !x.Equal(y) {
		t.Errorf("\nExpected: %+v\n     Got: %+v\n     Err: %v", expected, got, err)
	} else {
		fmt.Printf("%+v\n", got)
	}
}

func TestConditions(t *testing.T) {
	ts := time.Now()
	in := nmea2k.ParsedMessage{
		nmea2k.RawMessage{&can.RawMessage{
			ts,
			3, 129026, 1, 255, 8, []byte{0x0, 0xF, 0xC2, 0x40, 0xD0, 0x89, 0x00, 0x00},
		}},
		68,
		nmea2k.DataMap{0: 0, 1: "True", 2: 0xF, 3: 123.4, 4: 5.3},
	}

	expected := update{
		source{129026, "/dev/actisense", ts, 1},
		[]value{{"~/navigation/courseOverGroundTrue", 123.4}},
	}

	got, err := mapdata.Delta(&in)

	x := MakeSet(got.Values)
	y := MakeSet(expected.Values)

	if !x.Equal(y) {
		t.Errorf("\nExpected: %+v\n     Got: %+v\n     Err: %v", expected, got, err)
	} else {
		fmt.Printf("%+v\n", got)
	}
}

func TestRepeatingFields(t *testing.T) {
	ts := time.Now()
	in := nmea2k.ParsedMessage{
		nmea2k.RawMessage{&can.RawMessage{
			ts,
			3, 127503, 1, 255, 8, []byte{0x0, 0xF, 0xC2, 0x40, 0xD0, 0x89, 0x00, 0x00},
		}},
		85,
		nmea2k.DataMap{0: 0, 1: 3,
			2: "line1", 3: "Good", 5: 120.1, 6: 11, 7: 60, 8: 30, 9: 1321.1, 10: 1294.678, 11: 0.98,
			12: "line2", 13: "Bad Level", 15: 120.1, 16: 11, 17: 60, 18: 30, 19: 1321.1, 20: 1293.678, 21: 0.98,
			22: "line3", 23: "Bad Frequency", 25: 120.1, 26: 11, 27: 60, 28: 30, 29: 1321.1, 30: 1293.678, 31: 0.98,
		},
	}

	expected := update{
		source{127503, "/dev/actisense", ts, 1},
		[]value{
			{"~/electric/ac/0/numberOfLines", 3},
			{"~/electric/ac/0/line1/acceptability", "Good"},
			{"~/electric/ac/0/line2/acceptability", "Bad Level"},
			{"~/electric/ac/0/line3/acceptability", "Bad Frequency"},
		},
	}

	got, err := mapdata.Delta(&in)

	x := MakeSet(got.Values)
	y := MakeSet(expected.Values)

	if !x.Equal(y) {
		t.Errorf("\nExpected: %+v\n     Got: %+v\n     Err: %v", expected, got, err)
	} else {
		fmt.Printf("%+v\n", got)
	}
}
