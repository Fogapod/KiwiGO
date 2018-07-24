package config

import (
	"encoding/json"
	"os"
)

// Contains values read from config
type Config struct {
	Token        string `json:"token"`
	Prefix       string `json:"prefix"`
	LoggingLevel string `json:"logging_level"`
}

// Reads config file and updates it to match Config struct
func ReadConfig() (config *Config, err error) {
	f, err := os.OpenFile("config.json", os.O_RDWR, 0666)
	defer f.Close()

	if err != nil {
		return
	}

	b := make([]byte, 1024)
	var bytesRead int

	bytesRead, err = f.Read(b)

	if err != nil {
		return
	}

	if err = json.Unmarshal(b[:bytesRead], &config); err != nil {
		return
	}

	if b, err = json.MarshalIndent(config, "", "    "); err != nil {
		return
	}

	f.Seek(0, 0)
	_, err = f.Write(b)

	return
}
