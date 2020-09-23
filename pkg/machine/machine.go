package machine

import "errors"

// machine errors
var (
	ErrInvalidMachineState = errors.New("the machine is in invalid state")
)

type (
	// Machine represents a virtual machine
	Machine interface {
		// GetInfo gets information about the machine
		GetInfo() (*MachineInfo, error)

		// Start the machine
		Start() error

		// Stop the machine
		Stop() error
	}

	State string

	// MachineInfo holds information about a machine
	MachineInfo struct {
		Metadata struct {
			Name string
			ID   string
			State
		}
		Spec struct {
			OS     string
			Memory int
			VRAM   int
			CPUs   int
		}
	}
)

// machine states
const (
	PoweroffState State = "poweroff"
	RunningState        = "running"
	PausedState         = "paused"
	SavedState          = "saved"
	AbortedState        = "aborted"
	UnknownState        = "unknown"
)
