package events

// Type of event
type Type uint16

// Event an event
type Event struct {
	// Type the type of the event
	Type Type

	// Data the payload of the event
	Data interface{}
}

const (
	// TypeStartMachine the agent should start the vm when this event is received
	TypeStartMachine Type = iota

	// TypeStopMachine the agent should stop the vm when this event is received
	TypeStopMachine

	// TypeConnError connection error
	TypeConnError

	// UnknownEventError unknown event error
	UnknownEventError
)
