package virtualbox

import "sync"

type (
	// Options for the driver
	Options struct {
		// VirtualBoxManageCmdPath the path to the VBoxManage executable, by default
		// VBoxManage will be searched for in the $path
		VirtualBoxManageCmdPath string

		// StartWithGUI controls whether machines start headless mode
		// or with GUI, defaulf is headless
		StartWithGUI bool
	}

	driver struct {
		vboxManagePath string
		startWithGUI   bool
		mut            sync.Mutex
	}

	state string

	vm struct {
		d      *driver
		name   string
		uuid   string
		os     string
		state  state
		memory int
		cpus   int
		vram   int
	}
)

const (
	poweroffState state = "poweroff"
	runningState        = "running"
	pausedState         = "paused"
	savedState          = "saved"
	abortedState        = "aborted"
	unknownState        = "unknown"
)
