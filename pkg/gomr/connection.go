package gomr

import (
	"net"
)

type Connection struct {
	Hostname string
	Port     string
	Channel  string
	Nick     string
	Conn     net.Conn
}

func NewConnection(host, port, channel, nick string) (c *Connection, err error) {
	co := Connection{Hostname: host,
		Port:    port,
		Channel: channel,
		Nick:    nick}

	hostStr := co.Hostname + ":" + co.Port
	co.Conn, err = net.Dial("tcp", hostStr)
	if err != nil {
		return
	}
	co.Send("USER " + co.Nick + " 0 * " + co.Nick)
	co.Send("NICK " + co.Nick)
	co.Send("JOIN " + co.Channel)

	return &co, err
}

// Send the server a message
func (c *Connection) Send(text string) {
	//      log.Fprintf(conn, text)
	c.Conn.Write([]byte(text + "\n"))
}

// Identity can either be a channel or a nick
func (c *Connection) SendTo(identity, text string) {
	c.Conn.Write([]byte("PRIVMSG " + identity + " :" + text + "\n"))
}

// Send the configured channel a message
func (c *Connection) SendChan(text string) {
	c.Conn.Write([]byte("PRIVMSG " + c.Channel + " :" + text + "\n"))
}
