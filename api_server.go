package main

import (
	"github.com/timmathews/argo/nmea2k"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
)

type IndexEntry struct {
	Pgn     uint32
	Name    string
	Details string `json:"@Details"`
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<html><head><title>Pyxis API</title></head><body><h1>Pyxis API</h1></body></html>")
}

func MessagesIndex(w http.ResponseWriter, r *http.Request) {

	var pgns []IndexEntry

	category := r.URL.Query().Get("category")

	for i, pgn := range nmea2k.PgnList {
		if category != "" && !strings.EqualFold(category, pgn.Category) {
			continue
		}
		p := IndexEntry{}
		p.Pgn = pgn.Pgn
		p.Name = pgn.Description
		p.Details = fmt.Sprintf("/api/v1/messages/%d", i)
		pgns = append(pgns, p)
	}

	b, err := json.MarshalIndent(pgns, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Fprintf(w, string(b))
}

func MessageDetailsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, _ := vars["key"]
	tok := strings.Split(key, ".")
	pgn, _ := strconv.ParseInt(tok[0], 10, 32)
	var enc string
	var pgnDef interface{}
	var b []byte
	var err error

	if len(tok) > 1 {
		enc = tok[1]
	}

	if pgn >= int64(len(nmea2k.PgnList)) {
		pgnDef = map[string]interface{}{
			"Error": "Index out of range",
			"Index": pgn,
		}
	} else {
		pgnDef = nmea2k.PgnList[pgn]
	}

	if enc == "xml" {
		b, err = xml.Marshal(pgnDef)
	} else {
		b, err = json.MarshalIndent(pgnDef, "", "  ")
	}

	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Fprintf(w, string(b))
}

func ApiServer() {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/", HomeHandler)
	r.HandleFunc("/api/v1/messages", MessagesIndex)
	r.HandleFunc("/api/v1/messages/", MessagesIndex)
	r.HandleFunc("/api/v1/messages/{key}", MessageDetailsHandler)
	http.Handle("/api/v1/", r)
	http.ListenAndServe(":8082", nil)
}
