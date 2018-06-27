package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Token        string `json:"token"`
	Prefix       string `json:"prefix"`
	LoggingLevel int    `json:"logging_level"`
}

func readConfig() (Config, error) {
	var c Config

	raw, err := ioutil.ReadFile("./config.json")
	if err == nil {
		if err = json.Unmarshal(raw, &c); err == nil {
			return c, err
		}
	}

	return c, err
}
