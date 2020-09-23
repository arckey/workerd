package virtualbox

import (
	"bufio"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/arckey/workerd/pkg/drivers"
	"github.com/arckey/workerd/pkg/machine"
	"github.com/inconshreveable/log15"
)

// New creates a new virtualbox driver
func New(o *Options) (drivers.Driver, error) {
	var cmdPath string
	var err error
	if o != nil && o.VirtualBoxManageCmdPath != "" {
		cmdPath, err = exec.LookPath(o.VirtualBoxManageCmdPath)
	} else {
		cmdPath, err = exec.LookPath(vboxManageCmd)
	}

	if err != nil {
		return nil, ErrVBoxManageCommandNotFound
	}

	return &driver{
		vboxManagePath: cmdPath,
		mut:            sync.Mutex{},
	}, nil
}

// need to lock before running VBoxManage because it does not handle concurrent requests
// from the same process
func (d *driver) manage(args ...string) (string, string, error) {
	d.mut.Lock()
	defer d.mut.Unlock()

	// run the command
	cmd := exec.Command(d.vboxManagePath, args...)
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			log15.Error("command finished with error", "cmd", args, "error", string(exitErr.Stderr))
			return "", string(exitErr.Stderr), err
		}
		return "", "", err
	}
	return string(out), "", err
}

func (d *driver) GetMachineByName(name string) (machine.Machine, error) {
	out, _, err := d.manage("showvminfo", "--machinereadable", name)
	if err != nil {
		return nil, err
	}

	s := bufio.NewScanner(strings.NewReader(out))
	propMap := make(map[string]string)
	for s.Scan() {
		res := regexVMInfoLine.FindStringSubmatch(s.Text())
		if res == nil {
			continue
		}
		key := res[1]
		if key == "" {
			key = res[2]
		}
		val := res[3]
		if val == "" {
			val = res[4]
		}
		propMap[key] = val
	}

	if err = s.Err(); err != nil {
		return nil, err
	}

	vm := &vm{d: d}

	vm.name = propMap["name"]
	vm.uuid = propMap["UUID"]
	vm.os = propMap["ostype"]
	vm.memory, err = strconv.Atoi(propMap["memory"])
	vm.cpus, err = strconv.Atoi(propMap["cpus"])
	switch propMap["VMState"] {
	case "running":
		vm.state = runningState
	case "poweroff":
		vm.state = poweroffState
	case "paused":
		vm.state = pausedState
	case "aborted":
		vm.state = abortedState
	case "saved":
		vm.state = savedState
	default:
		vm.state = unknownState
	}

	return vm, nil
}

func (vm *vm) getTranslatedState() machine.State {
	switch vm.state {
	case runningState:
		return machine.RunningState
	case poweroffState:
		return machine.PoweroffState
	case pausedState:
		return machine.PausedState
	case abortedState:
		return machine.AbortedState
	case savedState:
		return machine.SavedState
	default:
		return machine.UnknownState
	}
}

func (vm *vm) unlockSession() error {
	_, _, err := vm.d.manage("startvm", vm.uuid, "--type", "emergencystop")
	return err
}

func (vm *vm) GetInfo() (*machine.MachineInfo, error) {

	return &machine.MachineInfo{
		Metadata: struct {
			Name string
			ID   string
			machine.State
		}{
			Name:  vm.name,
			ID:    vm.uuid,
			State: vm.getTranslatedState(),
		},
		Spec: struct {
			OS     string
			Memory int
			VRAM   int
			CPUs   int
		}{
			OS:     vm.os,
			Memory: vm.memory,
			VRAM:   vm.vram,
			CPUs:   vm.cpus,
		},
	}, nil
}

func (vm *vm) Start() error {
	var t string
	if vm.d.startWithGUI {
		t = "gui"
	} else {
		t = "headless"
	}

	_, stderr, err := vm.d.manage("startvm", vm.uuid, "--type", t)

	// checks if need to issue the unlock command
	if !strings.Contains(stderr, "is already locked by a session") {
		return err
	}
	if e := vm.unlockSession(); e != nil {
		return err // could not unlock machine session
	}

	_, _, err = vm.d.manage("startvm", vm.uuid, "--type", t) // try again

	return err
}

func (vm *vm) Stop() error {
	_, stderr, err := vm.d.manage("controlvm", vm.uuid, "poweroff")

	// checks if need to issue the unlock command
	if !strings.Contains(stderr, "is already locked by a session") {
		if strings.Contains(stderr, "is not currently running") {
			return machine.ErrInvalidMachineState
		}
		return err
	}
	if e := vm.unlockSession(); e != nil {
		return err // could not unlock machine session
	}
	_, _, err = vm.d.manage("controlvm", vm.uuid, "poweroff") // try again

	return err
}
