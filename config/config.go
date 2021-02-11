package config

import (
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

const (
	configFileName = "config.yaml"
)

type Config struct {
	Db struct {
		MaxOpenConnections int `yaml:"maxOpenConnections"`
	}
	Obfuscator struct {
		SliceSize         int   `yaml:"sliceSize"`
		DispersionPercent int64 `yaml:"dispersionPercent"`
	}
}

var config *Config

func GetConfig() Config {
	if config == nil {
		config = readConfigFromFile()
	}
	return *config
}

func readConfigFromFile() *Config {
	config := &Config{}

	configPath := getConfigPath()
	file, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)

	if err := decoder.Decode(&config); err != nil {
		panic(err)
	}
	return config
}

func getConfigPath() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return filepath.Join(wd, configFileName)
}
