package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Token        string `json:"token"`
	Prefix       string `json:"prefix"`
	LoggingLevel int    `json:"logging_level"`
}

func readConfig() (Config, error) {
	var c Config

	f, err := os.OpenFile("config.json", os.O_RDWR, 0666)
	defer f.Close()

	if err != nil {
		return c, err
	}

	b := make([]byte, 1024)
	var bytesRead int

	bytesRead, err = f.Read(b)

	if err != nil {
		return c, err
	}

	if err = json.Unmarshal(b[:bytesRead], &c); err != nil {
		return c, err
	}

	if b, err = json.MarshalIndent(c, "", "    "); err != nil {
		return c, err
	}

	f.Seek(0, 0)
	_, err = f.Write(b)

	return c, err
}
