package utils

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type TickersConfig struct {
	Tickers []string  `yaml:"tickers"`
	Job     JobConfig `yaml:"job"`
}

type JobConfig struct {
	Enabled                bool `yaml:"enabled"`
	TradingHoursOnly       bool `yaml:"trading_hours_only"`
	RefreshIntervalMinutes int  `yaml:"refresh_interval_minutes"`
}
type PostgresConfig struct {
	Mode          string `yaml:"mode"`
	DataDir       string `yaml:"data_dir"`
	BinaryPath    string `yaml:"binary_path"`
	StartupParams string `yaml:"startup_params"`
	LoadPath      string `yaml:"load_path"`
}

type ServerConfig struct {
	APIKey    string         `yaml:"api_key"`
	DBConnStr string         `yaml:"db_conn_str"`
	LogLevel  string         `yaml:"log_level"`
	Port      string         `yaml:"port"`
	Postgres  PostgresConfig `yaml:"postgres"`
}

func LoadServerConfig(path string) (*ServerConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	//log.Printf("Config file content:\n%s\n", string(data))
	var config ServerConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Printf("Error unmarshalling config file: %v\n", err)
		return nil, err
	}
	//fmt.Println(config)

	return &config, nil
}

func LoadTickersConfig(path string) (*TickersConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config TickersConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
func SaveServerConfig(path string, config *ServerConfig) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
