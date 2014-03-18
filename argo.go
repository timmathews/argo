/*
Argo collects data from NMEA-2000 and NMEA-0813[1] sensors and resends that
data over a LAN via ZeroMQ to be consumed by subscribers. Argo also calculates
additional data such as true wind speed, leeway, set, and drift. These values
are also sent over the network. Additionally, Argo can log data to a database
for later analysis.

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
	"argo/actisense"
  "argo/nmea2k"
	"flag"
	"fmt"
	zmq "github.com/alecthomas/gozmq"
	"github.com/schleibinger/sio"
	msgpack "github.com/vmihailenco/msgpack"
	"log"
	"strconv"
	"strings"
	"syscall"
	"time"
  "net/http"
)

// Timestamp format for printing
const layout = "2006-01-02-15:04:05.999"

/* The following startup command reverse engineered from Actisense NMEAreader.
 * It instructs the NGT1 to clear its PGN message TX list, thus it starts
 * sending all PGNs.
 */
var NGT_STARTUP_SEQ = []byte{0x11, 0x02, 0x00}

func main() {
	// Command line flags are defined here
	debug := flag.Bool("d", false, "Debug mode, extra logging information shown on stderr")
	verbose := flag.Bool("v", false, "Verbose mode, be chatty")
	help := flag.Bool("h", false, "This help message")
	pgn := flag.Int("pgn", 0, "Display only this PGN")
	src := flag.Int("source", 255, "Display PGNs from this source only")
	quiet := flag.Bool("q", false, "Don't display PGN data")
  addr := flag.String("addr", ":8081", "http service address")
	device := "/dev/ttyUSB0"

	flag.Parse()

	switch flag.NArg() {
	case 0:
		// Use default device
	case 1:
		device = flag.Arg(0)
	default:
		log.Fatalln("expected max 1 arg for the serial port device, default is %s", device)
	}

	if *help {
		flag.PrintDefaults()
		return
	}

	if *debug {
		log.Println("opening", device)
	}

	port, err := sio.Open(device, syscall.B115200)

	if err != nil {
		log.Fatalln("open: %s", err)
	}

	actisense.WriteMessage(port, actisense.NGT_MSG_SEND, NGT_STARTUP_SEQ)
	time.Sleep(2)

	rxbuf := []byte{0}
	rxch := make(chan byte)
	txch := make(chan nmea2k.ParsedMessage)

	go actisense.ReadNGT1(port, rxch, txch)

  // Start up the WebSockets hub
  go h.run()

	context, _ := zmq.NewContext()
	socket, _ := context.NewSocket(zmq.PUB)
	defer context.Close()
	defer socket.Close()
	socket.Bind("tcp://*:5555")

	go PgnDefServer(context)

  go WebSocketServer(addr)

  go ApiServer()

	// Print and transmit received messages
	go func() {
		for {
			res := <-txch

			if (*pgn == 0 || int(res.Header.Pgn) == *pgn) &&
				(*src == 255 || int(res.Header.Source) == *src) &&
				(*quiet == false) {
				if *debug {
					fmt.Println(res.Header.Print(*verbose))
				}
				fmt.Println(res.Print(*verbose))
			}

      bm := res.MsgPack()
      bj := res.JSON()

			socket.Send(bm, 0)
      h.broadcast <- bj
		}
	}()

	for {
		n, err := port.Read(rxbuf)
		if err != nil {
			log.Fatalln("Read error: ", err)
		}
		if n != len(rxbuf) {
			log.Fatalln("Short read ", n, len(rxbuf))
		}
		for _, b := range rxbuf {
			rxch <- b
		}
	}

}

func PgnDefServer(context zmq.Context) {
	socket, _ := context.NewSocket(zmq.REP)
	defer socket.Close()
	socket.Bind("tcp://*:5556")

	for {
		msg, _ := socket.Recv(0)

		tok := strings.Split(string(msg), ":")

		if tok[0] == "get" {
			pgn, _ := strconv.ParseUint(tok[1], 10, 32)
			pgnDef := nmea2k.PgnList[pgn]
			b, _ := msgpack.Marshal(pgnDef)
			socket.Send(b, 0)
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

