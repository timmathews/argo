package main

import "github.com/burntsushi/toml"

type tomlConfig struct {
	LogLevel   string `toml:"log_level"`
	MapFile    string
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

func ReadConfig(path string) (tomlConfig, error) {
	var config tomlConfig
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return config, err
	}

	return config, nil
}
