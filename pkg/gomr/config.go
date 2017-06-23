package gomr

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
	Source   string `yaml:"source"`

	// Dictionary Plugin
	WordnikAPIKey string `yaml:"wordnikapikey"`
}

type DbConfig struct {
	Hostname string `yaml:"hostname"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

// Set configuration variables from configuration file and environment
func GetConfiguration(path string) (c Config, d DbConfig) {
	// Bail if the configuration file isn't found
	//  Some configuration must be set via configuration file
	if _, err := os.Stat(path); os.IsNotExist(err) {
		glog.Fatalf("Configuration file %s not found.", path)
	}

	// Get inital configuration from yaml file
	c = getFileConfiguration(path)

	d = DbConfig{}
	// Set database configuration via env variables if they exist
	//  Environment variables take precedence over file configuration
	if e := os.Getenv("GOMR_DATABASE_SERVICE_HOST"); e != "" {
		d.Hostname = e
	}
	if e := os.Getenv("GOMR_DATABASE_SERVICE_PORT"); e != "" {
		d.Port = e
	}
	if e := os.Getenv("MYSQL_USER"); e != "" {
		d.Username = e
	}
	if e := os.Getenv("MYSQL_PASSWORD"); e != "" {
		d.Password = e
	}
	if e := os.Getenv("MYSQL_DATABASE"); e != "" {
		d.Name = e
	}

	return c, d
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
