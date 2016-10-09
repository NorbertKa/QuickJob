package config

import (
	"encoding/json"
	"errors"
	"os"
)

var (
	ErrOpenConfig   = errors.New("Couldn't open config file")
	ErrDecodeConfig = errors.New("Couldn't decode config file")
)

func readConfig(fileName string) (*os.File, error) {
	file, fileErr := os.Open(fileName)
	if fileErr != nil {
		return nil, fileErr
	}
	return file, nil
}

func OpenConfig(fileName string) (*Config, error) {
	configFile, err := readConfig(fileName)
	if err != nil {
		return nil, ErrOpenConfig
	}
	configDecoder := json.NewDecoder(configFile)
	config := Config{}
	err = configDecoder.Decode(&config)
	if err != nil {
		return nil, ErrDecodeConfig
	}
	return &config, nil
}
