package main

import ()

type ExamplePlugin struct {
}

// Register() should create any resources that the plugin may need
//   It should also check api endpoints to be sure the plugin will work
//   An error should be returned if the plugin should not be used.
func (e ExamplePlugin) Register() (err error) {
	return nil
}

// Parse() is the main function of a plugin. A string from the irc channel
//   will be provided as an argument, for the plugin to parse as it wishes.
// This should return a string that the bot should return.
func (e ExamplePlugin) Parse(user, channel, input string, conn *Connection) (err error) {
	if Match(input, "hello?") {
		conn.SendTo(channel, "Hello "+user+"!")
	}
	return nil
}

func (e ExamplePlugin) Help() (helpText string) {
	return "Example help text!"
}
