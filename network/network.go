package network

import (
	"fmt"
	"net"
)

/**
 * Broadcasts and receives on the network
 */
type Connection struct {
	udp  *net.UDPConn
	addr *net.UDPAddr
}

func NewConnection() (*Connection, error) {
	addr, err := net.ResolveUDPAddr("udp", "239.0.0.0:31337")
	if err != nil {
		return nil, fmt.Errorf("cannot resolve udp address or something %v", err)
	}

	udp, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, fmt.Errorf("cannot open udp port %v", err)
	}

	return &Connection{
		udp:  udp,
		addr: addr,
	}, nil
}

func (c *Connection) Close() error {
	return c.udp.Close()
}

func (c *Connection) Write(data []byte) (n int, err error) {
	return c.udp.WriteToUDP(data, c.addr)
}

func (c *Connection) Read() ([]byte, error) {
	buf := make([]byte, 1500)
	read, _, err := c.udp.ReadFromUDP(buf)
	if err != nil {
		return nil, err
	}
	return buf[:read], nil
}

/**
* encodes and decodes messages from a Connection
 */
type Protocol struct {
	inbox      chan string
	connection Connection
}
