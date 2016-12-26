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
	"bytes"
	"encoding/json"
	"github.com/burntsushi/toml"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"html/template"
	"io"
	"net/http"
	"strings"
)

type Uuid struct {
	Uuid []string `json:"uuid"`
}

type Vessel struct {
	Name         string
	Manufacturer string
	Model        string
	Year         int
	Registration string
	Mmsi         int
	Callsign     string
	Uuid         []string
}

type Connection struct {
	ListenOn string `json:"listen"`
	Port     int
}

type FormData struct {
	Vessel     Vessel
	Connection Connection
}

func uuidHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	uuid := Uuid{
		Uuid: strings.Split(uuid.NewV4().String(), "-"),
	}

	b, err := json.Marshal(uuid)

	if err != nil {
		log.Error(err)
		http.Error(w, "Could not generate UUID", 500)
	} else {
		io.WriteString(w, string(b))
	}
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var data FormData
		err := decoder.Decode(&data)
		if err != nil {
			log.Error(err)
			http.Error(w, "Could not parse data", 500) // What's the correct error here?
		} else {
			buf := new(bytes.Buffer)
			if err = toml.NewEncoder(buf).Encode(data); err != nil {
				log.Error(err)
			} else {
				log.Notice("\n", buf.String())
			}
		}
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, err := template.ParseFiles("templates/index.gtpl")
		if err == nil {
			t.Execute(w, nil)
		}
	}
}

func UiServer(addr *string, cmd chan CommandRequest) {
	r := mux.NewRouter()
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))
	r.HandleFunc("/admin/uuid", uuidHandler)
	r.HandleFunc("/admin", adminHandler)
	r.HandleFunc("/", indexHandler)
	http.Handle("/", r)
}
