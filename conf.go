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
	"fmt"
	"github.com/burntsushi/toml"
	"github.com/imdario/mergo"
)

type tomlConfig struct {
	LogLevel   string `toml:"log_level"`
	MapFile    string `toml:"map_file"`
	WebSockets webSocketsConfig
	Mqtt       mqttConfig
	Interfaces map[string]interfaceConfig
}

type webSocketsConfig struct {
	Enabled bool
	Port    int
}

type mqttConfig struct {
	Enabled bool
	UseSSL  bool
	Port    int
	Host    string
}

type interfaceConfig struct {
	Path  string
	Type  string
	Speed int
}

var defaultConfig = tomlConfig{
	LogLevel: "NONE",
	MapFile:  "map.xml",
	WebSockets: webSocketsConfig{
		Enabled: true,
		Port:    8082,
	},
	Mqtt: mqttConfig{
		Enabled: true,
		UseSSL:  true,
		Host:    "localhost",
		Port:    8883,
	},
	Interfaces: map[string]interfaceConfig{
		"actisense1": {
			Path:  "/dev/ttyUSB0",
			Type:  "actisense",
			Speed: 115200,
		},
	},
}

func ReadConfig(path string) (tomlConfig, error) {
	var config tomlConfig

	if _, err := toml.DecodeFile(path, &config); err != nil {
		return defaultConfig, err
	}

	fmt.Println("Config: ", config)

	fmt.Println("Default:", defaultConfig)

	if err := mergo.Merge(&config, defaultConfig); err != nil {
		return defaultConfig, err
	}

	fmt.Println("Merged: ", config)

	return config, nil
}
