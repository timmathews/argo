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
	"flag"
	"fmt"
	"strings"
)

type commandArgs struct {
	Help       bool
	Explain    bool
	Stats      bool
	Pgn        int
	Src        int
	LogLevel   string
	ConfigFile string
	MapFile    string
	DeviceType string
	DevicePath string
}

func GetCommandLineOptions() commandArgs {
	var args commandArgs

	flag.BoolVar(&args.Help, "help", false, "This help message")
	flag.BoolVar(&args.Explain, "explain", false, "Dump PGNs as JSON")
	flag.BoolVar(&args.Stats, "statistic", false, "Display live statistics")
	flag.IntVar(&args.Pgn, "pgn", 0, "Display only this PGN")
	flag.IntVar(&args.Src, "source", 255, "Display PGNs from this source only")
	flag.StringVar(&args.LogLevel, "log", "", "Set logging level: NONE, CRITICAL, ERROR, WARNING, NOTICE, INFO, DEBUG")
	flag.StringVar(&args.ConfigFile, "config", "argo.conf", "Path to config file")
	flag.StringVar(&args.MapFile, "map", "map.xml", "File to use for mapping between input and Signal K")
	flag.StringVar(&args.DeviceType, "device", "actisense", "Choose type of device: actisense, canusb, file")

	flag.Parse()

	args.LogLevel = strings.ToUpper(args.LogLevel)

	if flag.NArg() == 0 {
		if args.DeviceType == "file" {
			args.DevicePath = "sample.json"
		} else {
			args.DevicePath = "/dev/ttyUSB0"
		}
	} else {
		args.DevicePath = flag.Arg(0)
	}

	return args
}

func (c *commandArgs) PrintHelp() {
	fmt.Println("Argo Copyright (C) 2016 Tim Mathews <tim@signalk.org>")
	fmt.Println()
	flag.PrintDefaults()
}
