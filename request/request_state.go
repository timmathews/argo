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
)

type State int

const (
	Invalid State = iota
	Pending
	Completed
)

func (s State) String() string {
	return stateId[s]
}

var stateId = map[State]string{
	Invalid:   "INVALID",
	Pending:   "PENDING",
	Completed: "COMPLETED",
}

var stateName = map[string]State{
	"INVALID":   Invalid,
	"PENDING":   Pending,
	"COMPLETED": Completed,
}

func (s *State) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(stateId[*s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (s *State) UnmarshalJSON(b []byte) error {
	var v string
	err := json.Unmarshal(b, &v)

	if err != nil {
		return err
	}

	*s = stateName[v]
	return nil
}
