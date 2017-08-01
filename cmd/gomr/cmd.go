package main

import (
	"flag"

	"github.com/golang/glog"
	"github.com/tiwillia/gomr/pkg/gomr"
)

// This is where we start heh
func main() {
	// Base Configuration
	host := flag.String("host", "irc.freenode.net", "Hostname of the IRC server to connect to")
	port := flag.String("port", "6667", "Port of the IRC server to connect to")
	channel := flag.String("channel", "#test", "Name of the IRC channel to join")
	nick := flag.String("nick", "gomr", "Nickname of the IRC bot")
	password := flag.String("password", "", "IRC channel password (if applicable)")
	wordnikAPIKey := flag.String("wordnikapikey", "", "Wordnik API key for dictionary lookup support")
	source := flag.String("source", "https://github.com/tiwillia/gomr", "Source link for contribution recommendations")

	// Database configuration
	dbHost := flag.String("dbhost", "localhost", "Hostname of the mysql server to use")
	dbPort := flag.String("dbport", "3306", "Port of the mysql server to use")
	dbUsername := flag.String("dbusername", "gomr", "Username of the mysql user")
	dbPassword := flag.String("dbpassword", "", "Password of the mysql user")
	dbName := flag.String("dbname", "gomr", "Name of the mysql database")

	flag.Parse()
	glog.Infoln("Starting irc bot...")

	config := gomr.Config{
		Hostname:      *host,
		Port:          *port,
		Password:      *password,
		Channel:       *channel,
		Nick:          *nick,
		Source:        *source,
		WordnikAPIKey: *wordnikAPIKey,
	}

	dbConfig := gomr.DbConfig{
		Hostname: *dbHost,
		Port:     *dbPort,
		Username: *dbUsername,
		Password: *dbPassword,
		Name:     *dbName,
	}
	// Overwrite provided database configuration with environment variables
	dbConfig.GetEnv()

	gomrService, err := gomr.NewGomrService(&config, &dbConfig)
	if err != nil {
		glog.Fatalf("Unable to create Gomr service: %s", err)
	}

	err = gomrService.Run()
	if err != nil {
		glog.Fatalf("Error encountered running Gomr service: %s", err)
	}
}
