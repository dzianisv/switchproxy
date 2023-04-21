package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Rule struct {
	Domains []string `yaml:"domains"`
	Proxy   string   `yaml:"proxy"`
}

type Config struct {
	Rules  []Rule `yaml:"rules"`
	Listen string `yaml:"listen"`
}

func parseConfig(filename string) (Config, error) {
	var config Config

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	if config.Listen == "" {
		config.Listen = "127.0.0.1:8080"
	}

	return config, nil
}
