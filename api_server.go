package main

import (
  "fmt"
  "net/http"
  "github.com/gorilla/mux"
  "argo/nmea2k"
  "encoding/json"
  "strconv"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "<html><head><title>Pyxis API</title></head><body><h1>Pyxis API</h1></body></html>")
}

func MessagesIndex(w http.ResponseWriter, r *http.Request) {

  b, err := json.MarshalIndent(nmea2k.PgnList, "", "  ")
  if( err != nil ) {
    fmt.Println("error:", err)
  }
  fmt.Fprintf(w, string(b))
}

func MessageDetailsHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  pgn, _ := strconv.ParseInt(vars["key"], 10, 32)

  _, pgnDef := nmea2k.PgnList.First(uint32(pgn))

//  pgnDef := nmea2k.PgnList[pgn]
  b, err := json.MarshalIndent(pgnDef, "", "  ")
  if( err != nil) {
    fmt.Println("error:", err)
  }
  fmt.Fprintf(w, string(b))
}

func ApiServer() {
  r := mux.NewRouter()
  r.HandleFunc("/api/v1/", HomeHandler)
  r.HandleFunc("/api/v1/messages", MessagesIndex)
  r.HandleFunc("/api/v1/messages/{key}", MessageDetailsHandler)
  http.Handle("/api/v1/", r)
  http.ListenAndServe(":8082", nil)
}

