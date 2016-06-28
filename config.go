package main

import (
	"errors"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// Environment type database config
type Environment struct {
	Dialect    string `yaml:"dialect"`
	DataSource string `yaml:"datasource"`
	// Dir        string `yaml:"dir"`
	// TableName  string `yaml:"table"`
	// SchemaName string `yaml:"schema"`
}

// ReadConfig - reads configuration file
func ReadConfig(configFile string) (map[string]*Environment, error) {
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	config := make(map[string]*Environment)
	err = yaml.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// GetConfig from file for selected config
func GetConfig(configFile string, env string) (*Environment, error) {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return nil, nil
	}

	configMap, err := ReadConfig(configFile)
	if err != nil {
		return nil, err
	}

	e := configMap[env]
	if e == nil {
		return nil, errors.New("No environment: " + env)
	}

	if e.Dialect == "" {
		return nil, errors.New("No dialect specified")
	}

	if e.DataSource == "" {
		return nil, errors.New("No data source specified")
	}

	return e, nil
}
