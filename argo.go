/*
Argo collects data from NMEA-2000 and NMEA-0813[1] sensors and resends that
data over a LAN via an MQTT broker to be consumed by subscribers. Argo also
provides a WebSockets server Argo calculates additional data such as true wind
speed, leeway, set, and drift. These values are also sent over the network.
Additionally, Argo can log data to a database for later analysis.

Argo borrows heavily from the CANboat project which was written in C and is
copyright 2009-2012, Kees Verruijt, Harlingen, The Netherlands.

This file is part of Argo.

Argo is free software: you can redistribute it and/or modify it under the terms
of the GNU General Public License as published by the Free Software Foundation,
either version 3 of the License, or (at your option) any later version.

Argo is distributed in the hope that it will be useful, but WITHOUT ANY
WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A
PARTICULAR PURPOSE.  See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with
Argo.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/schleibinger/sio"
	"github.com/timmathews/argo/actisense"
	"github.com/timmathews/argo/can"
	"github.com/timmathews/argo/canusb"
	"github.com/timmathews/argo/nmea2k"
	"github.com/wsxiaoys/terminal"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"syscall"
	"time"
)

// Timestamp format for printing
const layout = "2006-01-02-15:04:05.999"

type UintSlice []uint32

func (p UintSlice) Len() int           { return len(p) }
func (p UintSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p UintSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func main() {
	// Command line flags are defined here
	debug := flag.Bool("d", false, "Debug mode, extra logging information shown on stderr")
	verbose := flag.Bool("v", false, "Verbose mode, be chatty")
	help := flag.Bool("h", false, "This help message")
	pgn := flag.Int("pgn", 0, "Display only this PGN")
	src := flag.Int("source", 255, "Display PGNs from this source only")
	quiet := flag.Bool("q", false, "Don't display PGN data")
	addr := flag.String("addr", ":8081", "http service address")
	stats := flag.Bool("s", false, "Display live statistics")
	dev_type := flag.String("dev", "actisense", "Choose type of device: actisense, canusb")
	no_server := flag.Bool("no-server", false, "Don't start Web Sockets or ZeroMQ")
	map_file := flag.String("map", "map.xml", "File to use for mapping between input and Signal K")
	device := "/dev/ttyUSB0"

	flag.Parse()

	if *dev_type != "canusb" && *dev_type != "actisense" {
		log.Fatalln("expected either canusb or actisense, got", *dev_type)
	}

	switch flag.NArg() {
	case 0:
		// Use default device
	case 1:
		device = flag.Arg(0)
	default:
		log.Fatalln("expected max 1 arg for the serial port device, default is", device)
	}

	if *help {
		flag.PrintDefaults()
		return
	}

	if *debug {
		log.Println("opening", device)
	}

	txch := make(chan nmea2k.ParsedMessage)
	cmdch := make(chan CommandRequest)

	statLog := make(map[uint32]uint64)
	var statPgns UintSlice

	data, err := ioutil.ReadFile(*map_file)
	if err != nil {
		log.Fatalln("could not read XML map file:", err, *map_file)
	}

	map_data := Mappings{}

	err = xml.Unmarshal(data, &map_data)
	if err != nil {
		log.Fatalln("could not parse XML map file:", err, *map_file)
	}

	port, err := sio.Open(device, syscall.B230400)

	if err != nil {
		log.Fatalln("open: %s", err)
	}

	var canport can.ReadWriter

	if !*no_server {
		// Start up the WebSockets hub
		go h.run()

		go WebSocketServer(addr)

		go ApiServer(cmdch)
	}

	// Print and transmit received messages
	go func() {
		for {
			res := <-txch

			if (*pgn == 0 || int(res.Header.Pgn) == *pgn) &&
				(*src == 255 || int(res.Header.Source) == *src) &&
				(*quiet == false) && (*stats == false) {
				if *debug {
					fmt.Println(res.Header.Print(*verbose))
				}
				fmt.Println(res.Print(*verbose))
			}

			if *stats {
				if _, ok := statLog[res.Header.Pgn]; ok {
					statLog[res.Header.Pgn]++
				} else {
					statLog[res.Header.Pgn] = 1
					statPgns = append(statPgns, res.Header.Pgn)
					sort.Sort(statPgns)
				}
				terminal.Stdout.Clear()
				for _, k := range statPgns {
					fmt.Println(k, "=>", statLog[k])
				}
			}

			if !*no_server {
				bj, err := map_data.Delta(&res)
				if err == nil {
					bytes, err := json.Marshal(bj)
					if err == nil {
						h.broadcast <- bytes
					}
				}
			}
		}
	}()

	// Set up hardware
	if *debug {
		fmt.Println(*dev_type)
	}
	if *dev_type == "canusb" {
		if *debug {
			fmt.Println("Adding Fast Packets")
		}
		for _, p := range nmea2k.PgnList {
			if p.Size > 8 {
				if *debug {
					log.Println("Adding PGN:", p.Pgn)
				}
				canusb.AddFastPacket(p.Pgn)
			}
		}

		if *debug {
			fmt.Println("Opening Channel")
		}
		canport, _ = canusb.OpenChannel(port, 221)
	} else {
		canport, _ = actisense.OpenChannel(port)
		time.Sleep(2)
	}

	// Handle command requests
	go func() {
		for {
			req := <-cmdch

			if req.RequestType == "iso" {
				b0 := (byte)(req.RequestedPgn) & 0xFF
				b1 := (byte)(req.RequestedPgn>>8) & 0xFF
				b2 := (byte)(req.RequestedPgn>>16) & 0xFF
				if canport != nil {
					canport.Write([]byte{0x03, 0x00, 0xEA, 0x00, 0xFF, 0x03, b0, b1, b2})
				} else {
					log.Println("canport is nil")
				}
			}
		}
	}()

	// Read from hardware
	if *dev_type == "canusb" {
		for {
			raw, err := canport.Read()
			if err == nil {
				txch <- *(nmea2k.ParsePacket(raw))
			}
		}
	} else {
		for {
			raw, err := canport.Read()
			if err == nil {
				txch <- *(nmea2k.ParsePacket(raw))
			}
		}
	}
}

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func WebSocketServer(addr *string) {
	http.HandleFunc("/ws/v1/", serveWs)
	err := http.ListenAndServe(*addr, Log(http.DefaultServeMux))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
