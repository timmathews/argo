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
