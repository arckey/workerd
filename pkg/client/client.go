package client

import "github.com/arckey/workerd/pkg/events"

// Client creates a connection to the specified host and allows subscribing for events
type Client interface {
	Connect() error
	Disconnect() error
	Chan() chan *events.Event
}

// Options the options to create the client
type Options struct {
	// HostAddr the host addr
	HostAddr string
}

// Type the type of client
type Type uint8

const (
	// TypeTCP client that uses tcp
	TypeTCP Type = iota
)

// New creates a new client
func New(t Type, o *Options) (Client, error) {
	// TODO add default options / check for nil

	switch t {
	case TypeTCP:
		return newTCPClient(o)
	default:
		return newTCPClient(o)
	}
}
