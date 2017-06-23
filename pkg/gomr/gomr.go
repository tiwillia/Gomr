package gomr

import (
	"bufio"
	"fmt"
	"reflect"
	"regexp"

	"github.com/go-gorp/gorp"
	"github.com/golang/glog"
)

type GomrService struct {
	Config  *Config
	Db      *gorp.DbMap
	Plugins []Plugin
}

func NewGomrService(config *Config, dbConfig *DbConfig) (*GomrService, error) {
	// Initiate database connection
	glog.Infoln("Getting database connection...")
	database, err := InitDB(dbConfig.Hostname, dbConfig.Port,
		dbConfig.Username, dbConfig.Password, dbConfig.Name)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to to database:", err)
	}

	// TODO allow plugins to be configurable somehow
	// Really I'd like to get rid of the word 'plugin' entirely since its a feature in go1.8beta now
	// That, or actually use the plugin feature. That would be neato
	var plugins []Plugin

	ex := ExamplePlugin{}
	plugins = append(plugins, ex)

	karma := KarmaPlugin{
		Db:   database,
		Nick: config.Nick,
	}
	plugins = append(plugins, karma)

	factoid := FactoidPlugin{
		// TODO this should be configurable
		Blacklist: []string{"why", "where", "who", "when", "how", "now"},
		Db:        database,
		Nick:      config.Nick,
	}
	plugins = append(plugins, factoid)

	dict := DictionaryPlugin{
		Nick:          config.Nick,
		WordnikAPIKey: config.WordnikAPIKey,
	}
	plugins = append(plugins, dict)

	service := &GomrService{
		Config:  config,
		Db:      database,
		Plugins: plugins,
	}

	err = service.RegisterPlugins()
	return service, err
}

func (s *GomrService) Run() error {
	// create a connection to the irc server and join channel
	conn, err := NewConnection(s.Config.Hostname, s.Config.Port, s.Config.Channel, s.Config.Nick)
	if err != nil {
		return fmt.Errorf("Unable to connect to", s.Config.Hostname, ":", s.Config.Port, err)
	}

	// Loop through the connection stream for the rest of forseeable time
	stream := bufio.NewReader(conn.Conn)
	for {
		line, err := stream.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				continue
			}
			glog.Infoln("ERROR: Failed to read from input stream:", err)
			return err // TODO eventually this should be continue, but I want errors to be very obvious for now
		}
		s.ParseLine(line, conn)
	}

	// Close the connection if we ever get here for some reason
	conn.Conn.Close()

	return nil
}

// All plugins should be registered here.
func (s *GomrService) RegisterPlugins() error {
	for _, p := range s.Plugins {
		err := p.Register()
		if err != nil {
			return fmt.Errorf("ERROR: Unable to register and enable plugin %s: %s\n", reflect.TypeOf(p), err)
		}
		glog.Infof("Successfully registered plugin %s", reflect.TypeOf(p))
	}
	return nil
}

// Main method to parse lines sent from the server
// Loops through each plugin in pluginList and runs the Parse() method from each
//   on the provided line
func (s *GomrService) ParseLine(line string, conn *Connection) {
	glog.Infoln(line)

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
		glog.Infoln("user:", user)
	}

	crgx = regexp.MustCompile(`\sPRIVMSG\s(\S+)\s`)
	cmatch := crgx.FindStringSubmatch(line)
	if cmatch != nil && len(cmatch) > 1 {
		channel = cmatch[1]
		if channel == s.Config.Nick {
			// This must be done to allow PRIVMSG's to users
			channel = user
		}
		glog.Infoln("channel:", channel)
	}

	mrgx = regexp.MustCompile(`\sPRIVMSG\s\S+\s:(.*)`)
	mmatch := mrgx.FindStringSubmatch(line)
	if mmatch != nil && len(mmatch) > 1 {
		msg = mmatch[1]
		glog.Infoln("message:", msg)
	}

	if msg != "" {
		// Check if the help command was sent
		if Match(msg, `(?i)`+s.Config.Nick+`[:,.]*\shelp`) {
			var helpText []string
			for _, plugin := range s.Plugins {
				texts := plugin.Help()
				helpText = append(helpText, texts...)
			}
			for _, text := range helpText {
				// For now, always send help text to the user in a private message
				//  It is likely the help text will get too big for a channel.
				conn.SendTo(user, text)
			}
			conn.SendTo(user, "Want to contribute? Source: "+s.Config.Source)
			if channel != user {
				conn.SendTo(channel, user+", help information sent via private message")
			}
			return
		}

		for _, p := range s.Plugins {
			err := p.Parse(user, channel, msg, conn)
			if err != nil {
				glog.Infoln("ERROR in plugin", reflect.TypeOf(p), ":", err)
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
		glog.Fatalln("Could not find host to ping in received ping string: ", line)
	}

	conn.Send("PONG " + pongHost)
	glog.Infoln("PONG " + pongHost)
}
