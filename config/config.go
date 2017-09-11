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

package config

import (
	"bytes"
	"github.com/burntsushi/toml"
	"github.com/imdario/mergo"
	"os"
)

type TomlConfig struct {
	LogLevel   string `toml:"log_level"`
	MapFile    string `toml:"map_file"`
	AssetPath  string `timp:"asset_path"`
	WebSockets webSocketsConfig
	Mqtt       mqttConfig
	Interfaces map[string]interfaceConfig
	Vessel     vesselConfig
}

type webSocketsConfig struct {
	Disabled bool
	Port     int
}

type mqttConfig struct {
	Disabled     bool
	UseCleartext bool
	Port         int
	Host         string
	ClientId     string
	Username     string
	Password     string
}

type interfaceConfig struct {
	Path  string
	Type  string
	Speed int
}

type vesselConfig struct {
	Name         string
	Make         string
	Model        string
	Year         int
	Mmsi         int
	Callsign     string
	Registration string
	Uuid         string
}

var defaultConfig = TomlConfig{
	LogLevel:  "INFO",
	MapFile:   "map.xml",
	AssetPath: "./assets",
	WebSockets: webSocketsConfig{
		Disabled: false,
		Port:     8082,
	},
	Mqtt: mqttConfig{
		Disabled:     false,
		UseCleartext: false,
		Host:         "localhost",
		Port:         8883,
		ClientId:     "argo",
		Username:     "signalk",
		Password:     "signalk",
	},
	Interfaces: map[string]interfaceConfig{
		"actisense1": {
			Path:  "/dev/ttyUSB0",
			Type:  "actisense",
			Speed: 115200,
		},
	},
}

func ReadConfig(path string) (TomlConfig, error) {
	var config TomlConfig

	if _, err := toml.DecodeFile(path, &config); err != nil {
		return defaultConfig, err
	}

	if err := mergo.Merge(&config, defaultConfig); err != nil {
		return defaultConfig, err
	}

	return config, nil
}

func WriteConfig(path string, config TomlConfig) error {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(config); err != nil {
		return err
	} else {
		if f, err := os.Create(path); err != nil {
			return err
		} else {
			f.Write(buf.Bytes())
			f.Close()
		}
	}

	return nil
}
