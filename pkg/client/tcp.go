package client

import (
	"net"

	"github.com/arckey/workerd/pkg/events"
)

type tcpClient struct {
	hostAddr *net.TCPAddr
	conn     *net.TCPConn
	c        chan *events.Event
	stopChan chan struct{}
}

func newTCPClient(o *Options) (Client, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", o.HostAddr)
	if err != nil {
		return nil, err
	}

	return &tcpClient{
		hostAddr: tcpAddr,
		conn:     nil,
		c:        make(chan *events.Event),
		stopChan: make(chan struct{}),
	}, nil
}

func (c *tcpClient) Connect() error {
	conn, err := net.DialTCP("tcp", nil, c.hostAddr)
	if err != nil {
		return err
	}
	c.conn = conn
	go c.handleConnection()
	return nil
}

func (c *tcpClient) Disconnect() error {
	c.stopChan <- struct{}{} // signal stop
	<-c.stopChan             // wait for confirmation
	return c.conn.Close()
}

func (c *tcpClient) Chan() chan *events.Event {
	return c.c
}

func (c *tcpClient) handleConnection() {
	buf := make([]byte, 64)
Loop:
	for {
		select {
		case <-c.stopChan:
			break Loop
		default:
		}

		size, err := c.conn.Read(buf)
		if err != nil {
			c.c <- &events.Event{
				Type: events.TypeConnError,
				Data: err,
			}
		}

		// parse the data
		for i := 0; i < size; i++ {
			ev := buf[i]
			switch ev {
			case 'u': // up
				c.c <- &events.Event{
					Type: events.TypeStartMachine,
					Data: nil,
				}
			case 'd': // down
				c.c <- &events.Event{
					Type: events.TypeStopMachine,
					Data: nil,
				}
			case '\n': // ignore
			default: // unknown event
				c.c <- &events.Event{
					Type: events.UnknownEventError,
					Data: ev,
				}
			}
		}
	}

	c.stopChan <- struct{}{} // signal finished
}
