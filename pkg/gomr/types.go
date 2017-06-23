package gomr

import ()

// All plugins should implement this interface
type Plugin interface {
	Register() error
	Parse(string, string, string, *Connection) error // Parse(user, channel, msg, connection)
	Help() []string
}
