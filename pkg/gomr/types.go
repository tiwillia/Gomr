package gomr

import (
	"os"
	"strings"
)

// All plugins should implement this interface
type Plugin interface {
	Register() error
	Parse(string, string, string, *Connection) error // Parse(user, channel, msg, connection)
	Help() []string
}

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

func (d *DbConfig) GetEnv() {
	// Set database configuration via env variables if they exist
	//  Environment variables take precedence
	dbService := strings.ToUpper(os.Getenv("DATABASE_SERVICE_NAME"))
	if e := os.Getenv(dbService + "_SERVICE_HOST"); e != "" {
		d.Hostname = e
	}
	if e := os.Getenv(dbService + "_SERVICE_PORT"); e != "" {
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
}
