package main

import (
	"bufio"
	"log"
)

// All plugins should implement this interface
type Plugin interface {
	Register() error
	Parse(string, *Connection) error
	Help() string
}

var (
	database   Database
	config     Config
	pluginList []Plugin
)

// This is where we start heh
func main() {
	log.Println("Starting irc bot...")

	// Read configuration
	configPath := "./gomr.yaml"
	config = GetConfiguration(configPath)

	log.Println("Getting database connection...")
	database, err := InitDB(config.Db.Hostname, config.Db.Port,
		config.Db.Username, config.Db.Password, config.Db.Name)
	if err != nil {
		log.Panicln("ERROR: Unable to connect to to database:", err)
	}

	// TODO remove this, only here so the damn thing builds
	log.Println("remove me", database)

	// Register all plugins
	log.Println("Registering plugins...")
	registerPlugins()

	// create a connection to the irc server and join channel
	var conn *Connection
	conn, err = NewConnection(config.Hostname, config.Port, config.Channel, config.Nick)
	if err != nil {
		log.Panicln("Unable to connect to", config.Hostname, ":", config.Port, err)
	}

	// Loop through the connection stream for the rest of forseeable time
	stream := bufio.NewReader(conn.Conn)
	for {
		line, err := stream.ReadString('\n')
		if err != nil {
			log.Println("Oh shit, an error occured: ")
			log.Println(err)
			return
		}
		parseLine(line, conn)
	}

	// Close the connection if we ever get here for some reason
	conn.Conn.Close()
}

// All plugins should be registered here.
func registerPlugins() {
	exPlugin := ExamplePlugin{}
	err := exPlugin.Register()
	if err != nil {
		log.Println("Unable to register example plugin, skipping plugin.")
	} else {
		pluginList = append(pluginList, exPlugin)
	}

	karmaPlugin := KarmaPlugin{}
	err = karmaPlugin.Register()
	if err != nil {
		log.Println("Unable to register karma plugin, skipping plugin.")
	} else {
		pluginList = append(pluginList, karmaPlugin)
	}
}

// Main method to parse lines sent from the server
// Loops through each plugin in pluginList and runs the Parse() method from each
//   on the provided line
func parseLine(line string, conn *Connection) {
	log.Printf(line)

	// If a PING is received from the server, respond to avoid being disconnected
	if Match(line, "PING :"+config.Hostname+"$") {
		respondToPing(line, conn)
	} else {
		message := MatchAndPull(line, "PRIVMSG", `PRIVMSG `+config.Channel+` :(.+)\n`)
		if message != "" {
			for _, plugin := range pluginList {
				plugin.Parse(line, conn)
			}
		}
	}
}

// Respond to pings from the irc server to keep the server alive
func respondToPing(ping string, conn *Connection) {
	conn.Send("PONG " + config.Hostname)
	log.Println("PONG " + config.Hostname)
}
