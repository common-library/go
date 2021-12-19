// Package socket provides a socket clent interface.
package socket

import (
	"net"
)

// Client is object that provides client infomation.
type Client struct {
	connnetion net.Conn
}

// Connect is connect to the address.
//
// ex) client.Connect("127.0.0.1:22222")
func (client *Client) Connect(network string, address string) error {
	connnetion, err := net.Dial(network, address)
	if err != nil {
		return err
	}

	client.connnetion = connnetion

	return nil
}

// Read is read data from connection.
//
// ex) readData, err := client.Read(1024)
func (client *Client) Read(recvSize int) (string, error) {
	buffer := make([]byte, recvSize)

	recvLen, err := client.connnetion.Read(buffer)
	if err != nil {
		return "", err
	}

	return string(buffer[:recvLen]), nil
}

// Write is write data to connection.
//
// ex) serverClient.Write("example")
func (client *Client) Write(data string) (int, error) {
	return client.connnetion.Write([]byte(data))
}

// Close is close the connection.
//
// ex) client.Close()
func (client *Client) Close() error {
	if client.connnetion == nil {
		return nil
	}

	err := client.connnetion.Close()
	client.connnetion = nil

	return err
}
