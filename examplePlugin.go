package main

import (
	"errors"
)

type ExamplePlugin struct {
}

// If a table in the database is needed, a table should be defined here with a struct.
//   The table must then be added to the database in gorp.go

// Register() should create any resources that the plugin may need
//   It should also, for example, check api endpoints to be sure the plugin will work
//   An error should be returned if the plugin should not be used.
func (e ExamplePlugin) Register() (err error) {
	return nil
}

// Parse() is the main function of a plugin. A string from the irc channel
//   will be provided as an argument, for the plugin to parse as it wishes.
// This should return a string that the bot should return.
func (e ExamplePlugin) Parse(user, channel, input string, conn *Connection) (err error) {

	// Check out the utils.go file for ease-of-use functions like Match()
	if Match(input, "^hello?") {
		conn.SendTo(channel, "Hello "+user+"!")
	}

	// What if an error occurs?
	// use the errors package to create a new error and return it:
	if Match(input, "^test error") {
		return errors.New("This is an example error that will be logged!")
		// If additional logging is necessary, you can import the "log" class and log yourself
	}

	return
}

// Help should return a slice of strings, each will be sent to the user requesting help
//  on seperate lines.
func (e ExamplePlugin) Help() (texts []string) {
	texts = append(texts, "Example help text!")
	return texts
}
