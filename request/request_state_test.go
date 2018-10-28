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

package requestState

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestMarshalValidState(t *testing.T) {
	s := []State{Pending, Completed}

	r, err := json.Marshal(s)

	if err != nil {
		t.Error(err)
	}

	if bytes.Compare(r, []byte(`["PENDING","COMPLETED"]`)) != 0 {
		t.Errorf(`Expected %v, got %s`, s[0], r)
	}
}

func TestUnmarshalValidState(t *testing.T) {
	s := []byte(`["PENDING","COMPLETED"]`)

	var r []State
	err := json.Unmarshal(s, &r)

	if err != nil {
		t.Error(err)
	}

	if len(r) != 2 || r[0] != Pending || r[1] != Completed {
		t.Error(`Expected [Pending, Completed], got`, r)
	}
}

// Everything but PENDING and COMPLETED should unmarshal to INVALID
func TestUnmarshalInvalidState(t *testing.T) {
	s := []byte(`["FOO"]`)

	var r []State
	err := json.Unmarshal(s, &r)

	if err != nil {
		t.Error(err)
	}

	if r[0] != Invalid {
		t.Error("Expected [], got", r[0])
	}
}
