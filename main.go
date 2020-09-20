package main

import (
	"os"

	"github.com/arckey/workerd/pkg/client"
	"github.com/arckey/workerd/pkg/config" // loads configuration
	"github.com/arckey/workerd/pkg/drivers"
	"github.com/arckey/workerd/pkg/events"
	"github.com/arckey/workerd/pkg/machine"
	log "github.com/inconshreveable/log15"
)

const (
	clientType = client.TypeTCP
	driverType = drivers.Virtualbox
)

func main() {
	// initialize client
	log.Debug("creating client", "type", clientType)
	c, err := client.New(client.TypeTCP, &client.Options{
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
	log.Debug("initializing driver driver", "type", driverType)
	d, err := drivers.New(drivers.Virtualbox, nil)
	if err != nil {
		log.Crit("failed to initialize driver", "type", driverType, "err", err)
		os.Exit(1)
	}

	log.Debug("validating virtual-machine exists:", "name", config.Config.MachineName)
	m := machine.GetByName(config.Config.MachineName, d)
	if _, err := m.GetInfo(); err != nil {
		log.Crit("could not validate virtual-machine exists", "type", driverType, "name", config.Config.MachineName, "err", err)
		os.Exit(1)
	}

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
			if err := m.Start(); err != nil {
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
