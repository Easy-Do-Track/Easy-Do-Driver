package main

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Server  ServerConfig  `toml:"server"`
	Tracker TrackerConfig `toml:"tracker"`
}

type ServerConfig struct {
	Address string `toml:"address"`
}

type Mapping struct {
	Name     string `toml:"name"`
	Rotation string `toml:"rotation"`
	Offset   Euler  `toml:"offset"`
	Multi    Euler  `toml:"multi"`
}

type TrackerConfig struct {
	Address  string             `toml:"address"`
	Mappings map[string]Mapping `toml:"mappings"`
}

func ConfigFromFile(filename string) (Config, error) {
	var conf Config

	_, err := toml.DecodeFile(filename, &conf)
	return conf, err
}
