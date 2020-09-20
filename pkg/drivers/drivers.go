package drivers

import (
	"errors"
)

// Driver a vm interface driver
type Driver interface {
	GetMachineInfo(hostname string) (*MachineInfo, error)
	StartMachine(hostname string) error
	StopMachine(hostname string) error
	RestartMachine(hostname string) error
}

// Type the driver's type
type Type string

// Options to initialize the driver with
type Options struct{}

const (
	// Virtualbox vb driver
	Virtualbox Type = "virtualbox"
)

type driverInitializer func(*Options) (Driver, error)

var driversMap = map[Type]driverInitializer{
	Virtualbox: newVirtualboxDriver,
}

var (
	// ErrDriverNotFound the specified driver type does not exist
	ErrDriverNotFound = errors.New("a driver with the specified type cound not be found")

	// ErrMachineNotFound the machine cannot be found
	ErrMachineNotFound = errors.New("machine cannot not be found")
)

// MachineInfo contains information about a machine
type MachineInfo struct {
	// Name the machine's name
	Name string

	// Spec the machine's specifications
	Spec struct{} // TODO: add struct here
}

// New creates a new driver of the specified type with the specified options
func New(typ Type, o *Options) (Driver, error) {
	init, ok := driversMap[typ]
	if !ok {
		return nil, ErrDriverNotFound
	}
	return init(o)
}
