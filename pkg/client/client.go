package client

import (
	"errors"

	"github.com/arckey/workerd/pkg/events"
)

var (
	errClientDoesNotExists = errors.New("a client of the specified type does not exist")
)

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
	// TCP client that uses tcp
	TCP Type = iota

	// Signals client that is controlled by os signals
	Signals
)

var clientsMap = map[Type]func(*Options) (Client, error){
	TCP:     newTCPClient,
	Signals: newSignalsClient,
}

// New creates a new client
func New(t Type, o *Options) (Client, error) {
	clientConstructor, exists := clientsMap[t]
	if !exists {
		return nil, errClientDoesNotExists
	}

	return clientConstructor(o)
}
