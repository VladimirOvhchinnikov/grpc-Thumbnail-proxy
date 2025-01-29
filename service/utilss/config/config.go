package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	LogFilePath       string              `json:"logFilePath"`
	Database          DatabaseConfig      `json:"database"`
	YoutubeClient     YouTubeClientConfig `json:"youtubeClient"`
	GRPCServerAddress string              `json:"grpcServerAddress"`
}

type DatabaseConfig struct {
	DataSourceName string `json:"dataSourceName"`
}

type YouTubeClientConfig struct {
	BaseURL      string `json:"baseUrl"`
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
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
