package main

import (
	"encoding/json"
	"github.com/timmathews/argo/nmea2k"
	"time"
)

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

func Delta(msg nmea2k.ParsedMessage) ([]byte, error) {

	delta := new(delta)

	delta.Context = "vessels.230099999" // TODO: Build this using MMSI from config file

	update := new(update)

	update.Source = source{
		msg.Header.Pgn,
		"/dev/actisense", // TODO: Get this from command line params or config file
		msg.Header.Timestamp,
		msg.Header.Source,
	}

	pgnFields := nmea2k.PgnList[msg.Index].FieldList

	for i, f := range pgnFields {
		if f.SignalkPath != "" && msg.Data[i] != nil {
			update.Values = append(update.Values, value{f.SignalkPath, msg.Data[i]})
		}
	}

	if update.Values != nil {
		delta.Updates = append(delta.Updates, *update)
	}

	if delta.Updates != nil {
		return json.Marshal(delta)
	} else {
		return nil, nil
	}
}
