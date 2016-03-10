package main

import (
	"bufio"
	"log"
	"reflect"
	"regexp"
)

// All plugins should implement this interface
type Plugin interface {
	Register() error
	Parse(string, string, string, *Connection) error // Parse(user, channel, msg, connection)
	Help() []string
}

var (
	config     Config
	pluginList []Plugin
)

// This is where we start heh
func main() {
	log.Println("Starting irc bot...")

	// Read configuration
	config = GetConfiguration()

	log.Println("Getting database connection...")
	err := InitDB(config.Db.Hostname, config.Db.Port,
		config.Db.Username, config.Db.Password, config.Db.Name)
	if err != nil {
		log.Panicln("ERROR: Unable to connect to to database:", err)
	}

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
			if err.Error() == "EOF" {
				continue
			}
			log.Println("Oh shit, an error occured:", err)
			return
		}
		parseLine(line, conn)
	}

	// Close the connection if we ever get here for some reason
	conn.Conn.Close()
}

// All plugins should be registered here.
func registerPlugins() {
	var plugins []Plugin

	exPlugin := ExamplePlugin{}
	plugins = append(plugins, exPlugin)

	karmaPlugin := KarmaPlugin{}
	plugins = append(plugins, karmaPlugin)

	factoidPlugin := FactoidPlugin{}
	plugins = append(plugins, factoidPlugin)

	dictPlugin := DictionaryPlugin{}
	plugins = append(plugins, dictPlugin)

	var err error
	for _, p := range plugins {
		err = p.Register()
		if err != nil {
			log.Println("ERROR: Unable to register and enable plugin", reflect.TypeOf(p), ":", err)
		} else {
			pluginList = append(pluginList, p)
		}
	}
}

// Main method to parse lines sent from the server
// Loops through each plugin in pluginList and runs the Parse() method from each
//   on the provided line
func parseLine(line string, conn *Connection) {
	log.Printf(line)

	// If a PING is received from the server, respond to avoid being disconnected
	if Match(line, "^PING :") {
		respondToPing(line, conn)
		return
	}

	// In this block, we use regex to determin the user who sent the message,
	//   the channel the message was sent on, and the message itself.
	// Example lines from server:
	// 2016/02/22 13:37:58 :tim!~tim@example.com PRIVMSG #test11123 :This is a test string
	// 2016/02/22 13:38:11 :tim!~tim@example.com NICK :timbo
	// 2016/02/22 13:38:13 :timbo!~tim@example.com PRIVMSG #test11123 :this is another test string
	var user, channel, msg string
	var urgx, crgx, mrgx *regexp.Regexp

	urgx = regexp.MustCompile(`:(\S+)!~`)
	umatch := urgx.FindStringSubmatch(line)
	if umatch != nil && len(umatch) > 1 {
		user = umatch[1]
		log.Println("user:", user)
	}

	crgx = regexp.MustCompile(`\sPRIVMSG\s(\S+)\s`)
	cmatch := crgx.FindStringSubmatch(line)
	if cmatch != nil && len(cmatch) > 1 {
		channel = cmatch[1]
		if channel == config.Nick {
			// This must be done to allow PRIVMSG's to users
			channel = user
		}
		log.Println("channel:", channel)
	}

	mrgx = regexp.MustCompile(`\sPRIVMSG\s\S+\s:(.*)`)
	mmatch := mrgx.FindStringSubmatch(line)
	if mmatch != nil && len(mmatch) > 1 {
		msg = mmatch[1]
		log.Println("message:", msg)
	}

	if msg != "" {
		// Check if the help command was sent
		if Match(msg, `(?i)`+config.Nick+`[:,.]*\shelp`) {
			var helpText []string
			for _, plugin := range pluginList {
				texts := plugin.Help()
				helpText = append(helpText, texts...)
			}
			for _, text := range helpText {
				// For now, always send help text to the user in a private message
				//  It is likely the help text will get too big for a channel.
				conn.SendTo(user, text)
			}
			conn.SendTo(user, "Want to contribute? Source: "+config.Source)
			if channel != user {
				conn.SendTo(channel, user+", help information sent via private message")
			}
			return
		}

		for _, p := range pluginList {
			err := p.Parse(user, channel, msg, conn)
			if err != nil {
				log.Println("ERROR in plugin", reflect.TypeOf(p), ":", err)
			}
		}
	}
}

// Respond to pings from the irc server to keep the server alive
func respondToPing(line string, conn *Connection) {
	hrgx := regexp.MustCompile(`^PING :(\S+)`)
	hmatch := hrgx.FindStringSubmatch(line)
	var pongHost string
	if hmatch != nil && len(hmatch) > 1 {
		pongHost = hmatch[1]
	} else {
		panic("Could not find host to ping in received ping string: " + line)
	}

	conn.Send("PONG " + pongHost)
	log.Println("PONG " + pongHost)
}
