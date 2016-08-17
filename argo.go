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
	"bufio"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	mqtt "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"github.com/op/go-logging"
	"github.com/schleibinger/sio"
	"github.com/timmathews/argo/actisense"
	"github.com/timmathews/argo/can"
	"github.com/timmathews/argo/canusb"
	"github.com/timmathews/argo/nmea2k"
	"github.com/wsxiaoys/terminal"
	"io/ioutil"
	"os"
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

var log = logging.MustGetLogger("argo")
var log_format = logging.MustStringFormatter(
	"%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:-8s} %{id:03x}%{color:reset} %{message}",
)

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
	dev_type := flag.String("dev", "actisense", "Choose type of device: actisense, canusb, file")
	no_server := flag.Bool("no-server", false, "Don't start Web Sockets or ZeroMQ")
	map_file := flag.String("map", "map.xml", "File to use for mapping between input and Signal K")
	mqtt_server := flag.String("mqtt", "localhost", "Defaults to MQTT broker on localhost")
	config_file := flag.String("config", "argo.conf", "Path to config file")
	explain := flag.Bool("explain", false, "Dump PGNs as JSON")
	device := "/dev/ttyUSB0"

	flag.Parse()

	log_backend := logging.NewLogBackend(os.Stderr, "", 0)
	log_formatter := logging.NewBackendFormatter(log_backend, log_format)
	log_filter := logging.AddModuleLevel(log_formatter)

	if *debug {
		log_filter.SetLevel(logging.DEBUG, "")
	} else {
		log_filter.SetLevel(logging.WARNING, "")
	}

	logging.SetBackend(log_filter)

	if *help {
		flag.PrintDefaults()
		return
	}

	if *explain {
		bytes, err := json.MarshalIndent(nmea2k.PgnList, "", "  ")
		if err == nil {
			fmt.Println(bytes)
		} else {
			log.Fatal(err)
		}
		return
	}

	config, err := ReadConfig(*config_file)
	if err != nil {
		log.Fatalf("could not read config file %v: %v", *config_file, err)
	}

	if config.LogLevel == "NONE" {
		log_filter = logging.AddModuleLevel(logging.NewLogBackend(ioutil.Discard, "", 0))
	} else {
		lvl, err := logging.LogLevel(config.LogLevel)
		if err == nil {
			log_filter.SetLevel(lvl, "")
		} else {
			log.Warningf("Could not set log level to %v: %v", config.LogLevel, err)
		}
	}

	switch flag.NArg() {
	case 0:
		// Use default device
	case 1:
		device = flag.Arg(0)
	default:
		log.Fatal("expected max 1 arg for the serial port device, default is", device)
	}

	log.Debug("opening", device)

	var stat syscall.Stat_t
	var port *sio.Port
	err = syscall.Stat(device, &stat)

	if err != nil {
		log.Fatalf("failure to stat %v: %v", device, err)
	}

	if stat.Mode&syscall.S_IFMT == syscall.S_IFCHR {
		log.Debugf("%v is a serial port", device)
		port, err = sio.Open(device, syscall.B230400)
	} else {
		log.Debugf("%v is a file", device)
		*dev_type = "file"
	}

	if err != nil {
		log.Fatal("failure to", err)
	}

	txch := make(chan nmea2k.ParsedMessage)
	cmdch := make(chan CommandRequest)

	statLog := make(map[uint32]uint64)
	var statPgns UintSlice

	data, err := ioutil.ReadFile(*map_file)
	if err != nil {
		log.Fatalf("could not read XML map file: %v, %v", err, *map_file)
	}

	map_data := Mappings{}

	err = xml.Unmarshal(data, &map_data)
	if err != nil {
		log.Fatalf("could not parse XML map file: %v, %v", err, *map_file)
	}

	// Set up MQTT Client
	var mqtt_client *mqtt.Client
	if config.Mqtt.Enabled {
		mqtt_opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("ssl://%v:8883", *mqtt_server))
		mqtt_opts.SetClientID("argo") // TODO: This needs to be moved to config file
		mqtt_opts.SetUsername("signalk")
		mqtt_opts.SetPassword("signalk")
		mqtt_opts.SetTLSConfig(&tls.Config{MinVersion: tls.VersionTLS12})
		mqtt_client = mqtt.NewClient(mqtt_opts)
		if token := mqtt_client.Connect(); token.Wait() && token.Error() != nil {
			log.Fatal("MQTT:", token.Error())
		}
	}

	var canport can.ReadWriter
	var fileScanner *bufio.Scanner

	if !*no_server {
		// Start up the WebSockets hub
		go websocket_hub.run()

		go WebSocketServer(addr, log)

		go ApiServer(cmdch)
	}

	// Print and transmit received messages
	go func() {
		for {
			res := <-txch

			if (*pgn == 0 || int(res.Header.Pgn) == *pgn) &&
				(*src == 255 || int(res.Header.Source) == *src) &&
				(*quiet == false) && (*stats == false) {
				log.Debug(res.Header.Print(*verbose))
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

			bj, err := map_data.Delta(&res)
			if err == nil {
				bytes, err := json.Marshal(bj)
				if err == nil {
					if !*no_server {
						websocket_hub.broadcast <- bytes
					}

					if config.Mqtt.Enabled {
						mqtt_client.Publish("signalk/argo", 0, false, bytes) // TODO: This should be in config file
					}
				}
			}
		}
	}()

	// Set up hardware
	log.Debug("configuring", *dev_type)
	if *dev_type == "canusb" {
		log.Debug("adding Fast Packets")
		for _, p := range nmea2k.PgnList {
			if p.Size > 8 {
				log.Debug("adding PGN:", p.Pgn)
				canusb.AddFastPacket(p.Pgn)
			}
		}
		log.Debug("opening channel")
		canport, _ = canusb.OpenChannel(port, 221)
	} else if *dev_type == "actisense" {
		log.Debug("opening channel")
		canport, _ = actisense.OpenChannel(port)
		time.Sleep(2)
	} else if *dev_type == "file" {
		file, _ := os.Open(device)
		fileScanner = bufio.NewScanner(file)
	} else {
		log.Fatal("unknown device type %s. Expected one of: canusb, actisense, file", *dev_type)
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
					log.Warning("canport is nil")
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
	} else if *dev_type == "actisense" {
		for {
			raw, err := canport.Read()
			if err == nil {
				txch <- *(nmea2k.ParsePacket(raw))
			}
		}
	} else { // it's a file
		for fileScanner.Scan() {
			txt := fileScanner.Text()
			txch <- *(nmea2k.FromCanBoat(txt))
			time.Sleep(500 * time.Millisecond)
		}
	}
}
