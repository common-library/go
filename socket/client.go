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
func (this *Client) Connect(network, address string) error {
	connnetion, err := net.Dial(network, address)
	if err != nil {
		return err
	}
	this.connnetion = connnetion

	return nil
}

// Read is read data from connection.
//
// ex) readData, err := client.Read(1024)
func (this *Client) Read(recvSize int) (string, error) {
	if this.connnetion == nil {
		return "", errors.New("please call the Connect function first")
	}

	buffer := make([]byte, recvSize)

	recvLen, err := this.connnetion.Read(buffer)
	if err != nil {
		return "", err
	}

	return string(buffer[:recvLen]), nil
}

// Write is write data to connection.
//
// ex) writeLen, err := client.Write("example")
func (this *Client) Write(data string) (int, error) {
	if this.connnetion == nil {
		return -1, errors.New("please call the Connect function first")
	}

	return this.connnetion.Write([]byte(data))
}

// Close is close the connection.
//
// ex) err := client.Close()
func (this *Client) Close() error {
	if this.connnetion == nil {
		return nil
	}

	err := this.connnetion.Close()
	this.connnetion = nil

	return err
}

// GetRemoteAddr is get the local Addr
//
// ex) addr := client.GetLocalAddr()
func (this *Client) GetLocalAddr() net.Addr {
	if this.connnetion == nil {
		return nil
	}

	return this.connnetion.LocalAddr()
}

// GetRemoteAddr is get the remote Addr
//
// ex) addr := client.GetRemoteAddr()
func (this *Client) GetRemoteAddr() net.Addr {
	if this.connnetion == nil {
		return nil
	}

	return this.connnetion.RemoteAddr()
}
