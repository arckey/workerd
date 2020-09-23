package client

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/arckey/workerd/pkg/events"
)

type signalsClient struct {
	sc       chan os.Signal
	c        chan *events.Event
	stopChan chan struct{}
}

func newSignalsClient(o *Options) (Client, error) {
	return &signalsClient{
		sc:       make(chan os.Signal, 16),
		c:        make(chan *events.Event),
		stopChan: make(chan struct{}),
	}, nil
}

func (c *signalsClient) Connect() error {
	signal.Notify(c.sc) // listen to all signals
	go c.handleConnection()
	return nil
}

func (c *signalsClient) Disconnect() error {
	c.stopChan <- struct{}{} // signal stop
	<-c.stopChan             // wait for confirmation
	return nil
}

func (c *signalsClient) Chan() chan *events.Event {
	return c.c
}

func (c *signalsClient) handleConnection() {
Loop:
	for {
		select {
		case <-c.stopChan:
			break Loop
		case sig := <-c.sc:
			switch sig {
			case syscall.SIGUSR1:
				c.c <- &events.Event{
					Type: events.TypeStartMachine,
					Data: nil,
				}
			case syscall.SIGUSR2:
				c.c <- &events.Event{
					Type: events.TypeStopMachine,
					Data: nil,
				}
			default:
				// ignore...
			}
		}
	}

	c.stopChan <- struct{}{} // signal finished
}
