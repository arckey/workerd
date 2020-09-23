package drivers

import (
	"github.com/arckey/workerd/pkg/machine"
)

// Driver a vm interface driver
type Driver interface {
	GetMachineByName(name string) (machine.Machine, error)
}
