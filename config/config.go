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
	"strings"
)

type TomlConfig struct {
	LogLevel   string
	MapFile    string
	Server     serverConfig
	Mqtt       mqttConfig
	Interfaces map[string]interfaceConfig
	Vessel     VesselConfig
}

type serverConfig struct {
	EnableWebsockets bool
	Port             int
	ListenOn         string
	AssetPath        string
}

type mqttConfig struct {
	Enable   bool
	UseTls   bool
	Port     int
	Host     string
	ClientId string
	Username string
	Password string
	Channel  string
}

type interfaceConfig struct {
	Path  string
	Type  string
	Speed int
}

type VesselConfig struct {
	Name         string
	Manufacturer string
	Model        string
	Year         int
	Mmsi         int
	Callsign     string
	Registration string
	Uuid         string
	Uuid0        string `toml:"-"`
	Uuid1        string `toml:"-"`
	Uuid2        string `toml:"-"`
	Uuid3        string `toml:"-"`
	Uuid4        string `toml:"-"`
}

var defaultConfig = TomlConfig{
	LogLevel: "INFO",
	MapFile:  "map.xml",
	Server: serverConfig{
		AssetPath:        "./assets",
		EnableWebsockets: true,
		Port:             8080,
	},
	Mqtt: mqttConfig{
		Enable:   false,
		UseTls:   true,
		Host:     "localhost",
		Port:     8883,
		ClientId: "argo",
		Username: "signalk",
		Password: "signalk",
		Channel:  "signalk/argo",
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

	u := strings.Split(config.Vessel.Uuid, "-")
	if len(u) == 5 {
		config.Vessel.Uuid0 = u[0]
		config.Vessel.Uuid1 = u[1]
		config.Vessel.Uuid2 = u[2]
		config.Vessel.Uuid3 = u[3]
		config.Vessel.Uuid4 = u[4]
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
