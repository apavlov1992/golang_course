package config

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type (
	Config struct {
		SourceUrl string `yaml:"SourceUrl"`
		DBFile    string `yaml:"DBFile"`
		IndexFile string `yaml:"IndexFile"`
	}
)

func NewConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("error reading YAML file: %v", err)
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("error unmarshalling YAML: %v", err)
	}
	return config, err
}
