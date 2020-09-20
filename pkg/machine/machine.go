package machine

import (
	"github.com/arckey/workerd/pkg/drivers"
)

// Machine a vm interface
type Machine struct {
	name   string
	driver drivers.Driver
}

// GetByName tries to get a vm using the given driver
func GetByName(name string, driver drivers.Driver) *Machine {
	return &Machine{
		name:   name,
		driver: driver,
	}
}

func (m *Machine) GetInfo() (*drivers.MachineInfo, error) {
	return m.driver.GetMachineInfo(m.name)
}

// Start starts the vm
func (m *Machine) Start() error {
	return m.driver.StartMachine(m.name)
}

// Stop stops the vm
func (m *Machine) Stop() error {
	return m.driver.StopMachine(m.name)
}

// Restart restarts the vm
func (m *Machine) Restart() error {
	return m.driver.RestartMachine(m.name)
}
