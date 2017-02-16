package main

import (
	"github.com/golang/glog"
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

// Set configuration variables from configuration file and environment
func GetConfiguration(path string) (c Config) {
	// Bail if the configuration file isn't found
	//  Some configuration must be set via configuration file
	if _, err := os.Stat(path); os.IsNotExist(err) {
		glog.Fatalf("Configuration file %s not found.", path)
	}

	// Get inital configuration from yaml file
	c = getFileConfiguration(path)

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
func getFileConfiguration(path string) (c Config) {
	confile, err := ioutil.ReadFile(path)
	if err != nil {
		glog.Fatalf("Failed to read configuration: %s", err.Error())
	}
	err = yaml.Unmarshal([]byte(confile), &c)
	if err != nil {
		glog.Fatalf("Failed to read configuration: %s", err.Error())
	}

	return c
}
