package utils

import (
	"encoding/json"
	"os"
)

type Config struct {
	LogFilePath       string `json:"logFilePath"`
	GrpcServerAddress string `json:"grpcServerAddress"`
}

func LoadConfig(filePath string) (*Config, error) {
	var config Config
	jsonFile, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonFile, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
