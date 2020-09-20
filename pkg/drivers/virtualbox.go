package drivers

import (
	"errors"

	vb "github.com/pitstopcloud/virtualbox-go"
)

type driver struct {
	vbox *vb.VBox
}

func newVirtualboxDriver(o *Options) (Driver, error) {
	vbox := vb.NewVBox(vb.Config{
		VirtualBoxPath: "/usr/local/bin/VBoxManage",
	})

	return &driver{
		vbox: vbox,
	}, nil
}

func (d *driver) GetMachineInfo(hostname string) (*MachineInfo, error) {
	vm, err := d.vbox.VMInfo(hostname)
	if err != nil {
		return nil, err
	}

	return &MachineInfo{
		Name: vm.Spec.Name,
		Spec: struct{}{},
	}, nil
}

func (d *driver) StartMachine(hostname string) error {
	_, err := d.vbox.Start(&vb.VirtualMachine{
		Spec: vb.VirtualMachineSpec{
			Name: hostname,
		},
	})

	return err
}

func (d *driver) StopMachine(hostname string) error {
	_, err := d.vbox.Stop(&vb.VirtualMachine{
		Spec: vb.VirtualMachineSpec{
			Name: hostname,
		},
	})

	return err
}

func (d *driver) RestartMachine(hostname string) error {
	return errors.New("not implemented")
}
