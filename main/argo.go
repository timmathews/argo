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
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jacobsa/go-serial/serial"
	"github.com/op/go-logging"
	"github.com/timmathews/argo/actisense"
	"github.com/timmathews/argo/canusb"
	"github.com/timmathews/argo/config"
	"github.com/timmathews/argo/nmea2k"
	"github.com/timmathews/argo/signalk"
	"github.com/wsxiaoys/terminal"
)

type StringSlice []string

func (p StringSlice) Len() int           { return len(p) }
func (p StringSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p StringSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

var log = logging.MustGetLogger("argo")
var logFormat = logging.MustStringFormatter(
	"%{color}%{time:15:04:05.000} ▶ %{level:4s} %{id:04d} %{message}%{color:reset}",
)

var sysconf config.TomlConfig
var statLog map[int]uint64

func main() {
	logBackend := logging.NewLogBackend(os.Stderr, "", 0)
	logFormatter := logging.NewBackendFormatter(logBackend, logFormat)
	logFilter := logging.AddModuleLevel(logFormatter)
	logging.SetBackend(logFilter)

	opts := GetCommandLineOptions()

	if opts.Help {
		opts.PrintHelp()
		return
	}

	if opts.Explain {
		bytes, err := json.MarshalIndent(nmea2k.PgnList, "", "  ")
		if err == nil {
			fmt.Println(string(bytes))
		} else {
			log.Fatal(err)
		}
		return
	}

	var err error
	sysconf, err = config.ReadConfig(opts.ConfigFile)
	if err != nil {
		log.Fatalf("could not read config file. %v", err)
	}

	if opts.LogLevel != "" {
		sysconf.LogLevel = opts.LogLevel
	}

	if sysconf.LogLevel == "NONE" {
		logFilter =
			logging.AddModuleLevel(logging.NewLogBackend(ioutil.Discard, "", 0))
		logging.SetBackend(logFilter)
	} else {
		requestedLogLevel := sysconf.LogLevel

		lvl, err := logging.LogLevel(requestedLogLevel)
		if err == nil {
			logFilter.SetLevel(lvl, "")
		} else {
			log.Warningf("Could not set log level to %v: %v", requestedLogLevel, err)
		}
	}

	log.Noticef("log level set to %v", logging.GetLevel(""))

	txch := make(chan nmea2k.ParsedMessage)
	cmdch := make(chan CommandRequest)

	statLog := make(map[string]uint64)
	var statPgns StringSlice

	mapData, err := signalk.ParseMappings(sysconf.MapFile)
	if err != nil {
		log.Fatalf("could not read XML map file %v: %v", sysconf.MapFile, err)
	}

	// Set up MQTT Client
	var mqttClient mqtt.Client
	if sysconf.Mqtt.Enable {
		mqttOpts := mqtt.NewClientOptions().AddBroker(
			fmt.Sprintf("ssl://%v:%v", sysconf.Mqtt.Host, sysconf.Mqtt.Port),
		)
		mqttOpts.SetClientID(sysconf.Mqtt.ClientId)
		mqttOpts.SetUsername(sysconf.Mqtt.Username)
		mqttOpts.SetPassword(sysconf.Mqtt.Password)
		if sysconf.Mqtt.UseTls {
			mqttOpts.SetTLSConfig(&tls.Config{MinVersion: tls.VersionTLS12})
		}
		mqttClient = mqtt.NewClient(mqttOpts)
		if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
			log.Fatal("MQTT:", token.Error())
		}
	}

	// Convert the port int to a string starting with :
	// TODO: specify interfaces to listen on
	addr := fmt.Sprintf(":%v", sysconf.Server.Port)

	if sysconf.Server.EnableWebsockets {
		// Start up the WebSockets hub
		go websocket_hub.run()

		go statistics_hub.run()

		go WebSocketServer(&addr, log)
	}

	go ApiServer(&addr, cmdch)
	go UiServer(&addr, cmdch)

	// Print and transmit received messages
	go func() {
		verbose := logging.GetLevel("") == logging.DEBUG

		for {
			res := <-txch

			if (opts.Pgn == 0 || int(res.Header.Pgn) == opts.Pgn) &&
				(opts.Src == 255 || int(res.Header.Source) == opts.Src) &&
				(opts.Dst == 255 || int(res.Header.Destination) == opts.Dst) &&
				!opts.Stats {
				log.Debug(res.Header.Print(verbose))
				log.Info(res.Print(verbose))
			}

			pgn := strconv.Itoa(int(res.Header.Pgn))

			if _, ok := statLog[pgn]; ok {
				statLog[pgn]++
			} else {
				statLog[pgn] = 1
				statPgns = append(statPgns, pgn)
				sort.Sort(statPgns)
			}

			if sysconf.Server.EnableWebsockets {
				if b, err := json.Marshal(statLog); err == nil {
					statistics_hub.broadcast <- b
				} else {
					log.Errorf("JSON.Marshal %v", err)
				}
			}

			if opts.Stats {
				terminal.Stdout.Clear()
				for _, k := range statPgns {
					fmt.Println(k, "=>", statLog[k])
				}
			}

			bj, err := mapData.Delta(&res)
			if err == nil {
				bytes, err := json.Marshal(bj)
				if err == nil {
					if sysconf.Server.EnableWebsockets {
						websocket_hub.broadcast <- bytes
					}

					if sysconf.Mqtt.Enable {
						mqttClient.Publish(sysconf.Mqtt.Channel, 0, false, bytes)
					}
				}
			}
		}
	}()

	for k, i := range sysconf.Interfaces {
		log.Noticef("opening %v at %v", k, i.Path)
		go processInterface(i, txch)
	}

	exitc := make(chan os.Signal, 1)
	signal.Notify(exitc, os.Interrupt, syscall.SIGTERM)

	sig := <-exitc

	log.Notice("cleaning up and exiting with %v", sig)
}

func processInterface(iface config.InterfaceConfig, txch chan nmea2k.ParsedMessage) {
	var stat syscall.Stat_t
	var port io.ReadWriteCloser

	err := syscall.Stat(iface.Path, &stat)
	if err != nil {
		log.Fatalf("failure to stat %v: %v", iface.Path, err)
	}

	if stat.Mode&syscall.S_IFMT == syscall.S_IFCHR {
		log.Debugf("%v is a serial port", iface.Path)

		options := serial.OpenOptions{
			PortName:        iface.Path,
			BaudRate:        iface.Speed,
			DataBits:        8,
			StopBits:        1,
			MinimumReadSize: 4,
		}

		port, err = serial.Open(options)

		if err != nil {
			log.Fatalf("error opening port:", err)
		}
	} else {
		log.Debugf("%v is a file", iface.Path)
	}

	// Set up hardware and start reading data
	log.Debug("configuring %v", iface.Type)

	if iface.Type == "canusb" {
		log.Debug("adding Fast Packets")

		for _, p := range nmea2k.PgnList {
			if p.Size > 8 {
				log.Debug("adding PGN: %d", p.Pgn)
				canusb.AddFastPacket(p.Pgn)
			}
		}

		// Read from hardware
		log.Debug("opening channel")

		canport, _ := canusb.OpenChannel(port, 221)

		for {
			raw, err := canport.Read()
			if err == nil {
				txch <- *(nmea2k.ParsePacket(raw))
			} else {
				log.Warning("canport:", err)
			}
		}
	} else if iface.Type == "actisense" {
		// Read from hardware
		log.Debug("opening channel")
		canport, _ := actisense.OpenChannel(port)
		time.Sleep(250)

		canport.GetOperatingMode()

		for {
			raw, err := canport.Read()
			if err == nil {
				txch <- *(nmea2k.ParsePacket(raw))
			} else {
				log.Warning("canport:", err)
			}
		}
	} else if iface.Type == "file" {
		// Read from file
		file, _ := os.Open(iface.Path)
		fileScanner := bufio.NewScanner(file)

		for fileScanner.Scan() {
			txt := fileScanner.Text()
			pgn, err := nmea2k.FromCanBoat(txt)
			if err == nil {
				txch <- *pgn
			} else {
				log.Warning("filescanner:", err)
			}
			time.Sleep(100 * time.Millisecond)
		}
	} else {
		log.Fatalf(
			"unknown device type %s. Expected one of: canusb, actisense, file",
			iface.Type,
		)
	}
}
