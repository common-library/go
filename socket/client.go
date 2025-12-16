// Package socket provides socket client and server implementations.
package socket

import (
	"errors"
	"net"
)

// Client is a struct that provides client related methods.
type Client struct {
	connnetion net.Conn
}

// Connect is connect to the address.
//
// ex) err := client.Connect("tcp", "127.0.0.1:10000")
func (c *Client) Connect(network, address string) error {
	connnetion, err := net.Dial(network, address)
	if err != nil {
		return err
	}
	c.connnetion = connnetion

	return nil
}

// Read is read data from connection.
//
// ex) readData, err := client.Read(1024)
func (c *Client) Read(recvSize int) (string, error) {
	if c.connnetion == nil {
		return "", errors.New("please call the Connect function first")
	}

	buffer := make([]byte, recvSize)

	recvLen, err := c.connnetion.Read(buffer)
	if err != nil {
		return "", err
	}

	return string(buffer[:recvLen]), nil
}

// Write is write data to connection.
//
// ex) writeLen, err := client.Write("example")
func (c *Client) Write(data string) (int, error) {
	if c.connnetion == nil {
		return -1, errors.New("please call the Connect function first")
	}

	return c.connnetion.Write([]byte(data))
}

// Close is close the connection.
//
// ex) err := client.Close()
func (c *Client) Close() error {
	if c.connnetion == nil {
		return nil
	}

	err := c.connnetion.Close()
	c.connnetion = nil

	return err
}

// GetRemoteAddr is get the local Addr
//
// ex) addr := client.GetLocalAddr()
func (c *Client) GetLocalAddr() net.Addr {
	if c.connnetion == nil {
		return nil
	}

	return c.connnetion.LocalAddr()
}

// GetRemoteAddr is get the remote Addr
//
// ex) addr := client.GetRemoteAddr()
func (c *Client) GetRemoteAddr() net.Addr {
	if c.connnetion == nil {
		return nil
	}

	return c.connnetion.RemoteAddr()
}
