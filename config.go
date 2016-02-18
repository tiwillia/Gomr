package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// This struct defines the yaml the configuration file must follow
type Config struct {
	Hostname string `yaml:"hostname"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	Channel  string `yaml:"channel"`
	Nick     string `yaml:"nick"`
	Db       struct {
		Hostname string `yaml:"hostname"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	}
}

// Set configuration variables from configuration file
func GetConfiguration(filePath string) (c Config) {
	confile, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic("Failed to read configuration: " + err.Error())
	}
	err = yaml.Unmarshal([]byte(confile), &c)
	if err != nil {
		panic("Failed to read configuration: " + err.Error())
	}

	return c
}
