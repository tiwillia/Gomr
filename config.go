package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
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
	Source string `yaml:"source"`

	// Dictionary Plugin
	WordnikApiKey string `yaml:"wordnikapikey"`
}

var (
	configFilePath = "./gomr.yaml"
)

// Set configuration variables from configuration file and environment
func GetConfiguration() (c Config) {
	// Bail if the configuration file isn't found
	//  Some configuration must be set via configuration file
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		panic("Configuration file " + configFilePath + " not found.")
	}

	// Get inital configuration from yaml file
	c = getFileConfiguration()

	// Set database configuration via env variables if they exist
	//  Environment variables take precedence over file configuration
	if e := os.Getenv("GOMR_DATABASE_SERVICE_HOST"); e != "" {
		c.Db.Hostname = e
	}
	if e := os.Getenv("GOMR_DATABASE_SERVICE_PORT"); e != "" {
		c.Db.Port = e
	}
	if e := os.Getenv("MYSQL_USER"); e != "" {
		c.Db.Username = e
	}
	if e := os.Getenv("MYSQL_PASSWORD"); e != "" {
		c.Db.Password = e
	}
	if e := os.Getenv("MYSQL_DATABASE"); e != "" {
		c.Db.Name = e
	}

	return c
}

// Get configuration from yaml file
func getFileConfiguration() (c Config) {
	confile, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		panic("Failed to read configuration: " + err.Error())
	}
	err = yaml.Unmarshal([]byte(confile), &c)
	if err != nil {
		panic("Failed to read configuration: " + err.Error())
	}

	return c
}
