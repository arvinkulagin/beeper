package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	WSAddress   string
	RESTAddress string
	Origin      string
	MessageType int
}

func Default() Config {
	return Config{
		WSAddress:   "localhost:8888",
		RESTAddress: "localhost:8889",
		Origin:      "",
		MessageType: 1,
	}
}

func FromFile(file string, cfg *Config) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, cfg)
	if err != nil {
		return err
	}
	return nil
}