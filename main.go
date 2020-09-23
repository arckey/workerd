package main

import (
	"os"

	"github.com/arckey/workerd/pkg/client"
	"github.com/arckey/workerd/pkg/config"
	"github.com/arckey/workerd/pkg/drivers/virtualbox"
	"github.com/arckey/workerd/pkg/events"
	log "github.com/inconshreveable/log15"
)

const (
	clientType = client.TCP
)

func main() {
	// initialize client
	log.Debug("creating client", "type", clientType, "pid", os.Getpid())
	c, err := client.New(client.TCP, &client.Options{
		HostAddr: config.Config.WMAddr,
	})
	if err != nil {
		log.Crit("could not create client", "err", err)
		os.Exit(1)
	}

	log.Debug("trying to connect...", "host", config.Config.WMAddr)
	if err = c.Connect(); err != nil {
		log.Crit("could not create client", "err", err)
		os.Exit(1)
	}
	log.Info("established connection to:", "host", config.Config.WMAddr)

	// initialize driver
	log.Debug("initializing driver driver")
	d, err := virtualbox.New(nil)
	if err != nil {
		log.Crit("failed to initialize driver", "err", err)
		os.Exit(1)
	}

	m, err := d.GetMachineByName(config.Config.MachineName)
	if err != nil {
		log.Crit("could not get machine", "machine", config.Config.MachineName, "err", err)
		os.Exit(1)
	}

	info, err := m.GetInfo()
	if err != nil {
		log.Crit("could not get machine info", "machine", config.Config.MachineName, "err", err)
		os.Exit(1)
	}
	log.Info("found virtual-machine", "name", info.Metadata.Name, "state", info.Metadata.State)

	log.Info("listening for events from worker-manager...", "host", config.Config.WMAddr)
	for {
		ev := <-c.Chan()
		switch ev.Type {
		case events.TypeStartMachine:
			log.Info("Starting machine", "machine", config.Config.MachineName)
			if err := m.Start(); err != nil {
				log.Error("failed to start machine", "machine", config.Config.MachineName, "err", err)
			}
		case events.TypeStopMachine:
			log.Info("Stopping machine", "machine", config.Config.MachineName)
			if err := m.Stop(); err != nil {
				log.Error("failed to stop machine", "machine", config.Config.MachineName, "err", err)
			}
		case events.UnknownEventError:
			typ, _ := ev.Data.(byte)
			log.Error("got unknown event type", "type", typ)
		case events.TypeConnError:
			err, _ := ev.Data.(error)
			log.Error("connection error", "err", err)
		}
	}
}
